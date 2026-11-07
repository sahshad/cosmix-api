const fs = require('node:fs');
const path = require('node:path');

const COLLECTION_SCHEMA =
  'https://schema.getpostman.com/json/collection/v2.1.0/collection.json';
const POSTMAN_API_URL = 'https://api.getpostman.com';
const ROOT_DIR = path.resolve(__dirname, '..');

function loadEnvFile(filePath) {
  if (!fs.existsSync(filePath)) {
    return;
  }

  const contents = fs.readFileSync(filePath, 'utf8');
  for (const line of contents.split(/\r?\n/)) {
    const match = line.match(/^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)?\s*$/);
    if (!match || Object.hasOwn(process.env, match[1])) {
      continue;
    }

    let value = match[2] ?? '';
    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    } else {
      value = value.replace(/\s+#.*$/, '').trim();
    }
    process.env[match[1]] = value;
  }
}

function resolveReference(document, reference) {
  if (!reference.startsWith('#/')) {
    throw new Error(`Unsupported external OpenAPI reference: ${reference}`);
  }

  return reference
    .slice(2)
    .split('/')
    .map((part) => part.replace(/~1/g, '/').replace(/~0/g, '~'))
    .reduce((current, part) => current?.[part], document);
}

function resolveSchema(schema, document, seen = new Set()) {
  if (!schema || typeof schema !== 'object') {
    return {};
  }

  if (schema.$ref) {
    if (seen.has(schema.$ref)) {
      return {};
    }
    const nextSeen = new Set(seen).add(schema.$ref);
    return resolveSchema(resolveReference(document, schema.$ref), document, nextSeen);
  }

  if (Array.isArray(schema.allOf)) {
    return schema.allOf.reduce(
      (combined, part) => ({
        ...combined,
        ...resolveSchema(part, document, seen),
        properties: {
          ...combined.properties,
          ...resolveSchema(part, document, seen).properties,
        },
        required: [
          ...(combined.required ?? []),
          ...(resolveSchema(part, document, seen).required ?? []),
        ],
      }),
      { ...schema, properties: {}, required: [] },
    );
  }

  return schema;
}

function sampleValue(schema, document, seen = new Set()) {
  if (schema?.$ref) {
    if (seen.has(schema.$ref)) {
      return {};
    }
    return sampleValue(
      resolveReference(document, schema.$ref),
      document,
      new Set(seen).add(schema.$ref),
    );
  }

  const resolved = resolveSchema(schema, document, seen);
  if (resolved.example !== undefined) {
    return resolved.example;
  }
  if (resolved.default !== undefined) {
    return resolved.default;
  }
  if (resolved.enum?.length) {
    return resolved.enum[0];
  }

  switch (resolved.type) {
    case 'object': {
      const required = new Set(resolved.required ?? []);
      const properties = resolved.properties ?? {};
      const keys = required.size ? [...required] : Object.keys(properties);
      return Object.fromEntries(
        keys
          .filter((key) => properties[key])
          .map((key) => [key, sampleValue(properties[key], document, seen)]),
      );
    }
    case 'array':
      return [sampleValue(resolved.items ?? {}, document, seen)];
    case 'integer':
    case 'number':
      return 0;
    case 'boolean':
      return false;
    case 'string':
      return resolved.format === 'date-time' ? '2026-01-01T00:00:00Z' : 'string';
    default:
      return {};
  }
}

function getPathParameters(pathTemplate) {
  return [...pathTemplate.matchAll(/\{([^}]+)\}/g)].map((match) => ({
    key: match[1],
    value: `{{${match[1]}}}`,
    description: 'Set this value before sending the request.',
  }));
}

function makeRequestUrl(baseUrl, pathTemplate, parameters) {
  const pathWithVariables = pathTemplate.replace(
    /\{([^}]+)\}/g,
    ':$1'
  );

  const url = new URL(
    `${baseUrl.replace(/\/+$/, '')}/${pathWithVariables.replace(/^\/+/, '')}`
  );

  const query = parameters
    .filter((parameter) => parameter.in === 'query')
    .map((parameter) => ({
      key: parameter.name,
      value: String(
        parameter.example ??
        parameter.schema?.example ??
        ''
      ),
      description: parameter.description,
      disabled: parameter.required !== true,
    }));

  const pathVariables = getPathParameters(pathTemplate).map((variable) => {
    const definition = parameters.find(
      (parameter) =>
        parameter.in === 'path' &&
        parameter.name === variable.key
    );

    return {
      ...variable,
      description: definition?.description ?? variable.description,
    };
  });

  return {
    raw: url.toString(),
    protocol: url.protocol.replace(':', ''),
    host: [url.hostname],
    ...(url.port ? { port: url.port } : {}),
    path: url.pathname.split('/').filter(Boolean),
    ...(query.length ? { query } : {}),
    ...(pathVariables.length ? { variable: pathVariables } : {}),
  };
}

function makeRequestBody(operation, document) {
  const content = operation.requestBody?.content ?? {};
  const contentType = Object.keys(content).find((type) => type.includes('json'));
  if (contentType) {
    const schema = content[contentType].schema ?? {};
    const example =
      content[contentType].example ??
      sampleValue(schema, document);
    return {
      mode: 'raw',
      raw: JSON.stringify(example, null, 2),
      options: { raw: { language: 'json' } },
    };
  }

  const formType = Object.keys(content).find(
    (type) =>
      type === 'application/x-www-form-urlencoded' || type === 'multipart/form-data',
  );
  if (!formType) {
    return undefined;
  }

  const schema = resolveSchema(content[formType].schema ?? {}, document);
  const fields = Object.entries(schema.properties ?? {}).map(([key, property]) => ({
    key,
    value: String(sampleValue(property, document)),
    type: 'text',
    disabled: !(schema.required ?? []).includes(key),
  }));
  return {
    mode: formType === 'multipart/form-data' ? 'formdata' : 'urlencoded',
    [formType === 'multipart/form-data' ? 'formdata' : 'urlencoded']: fields,
  };
}

function makePostmanItem(document, baseUrl, pathTemplate, method, operation, pathParameters) {
  const parameters = [...(pathParameters ?? []), ...(operation.parameters ?? [])];
  const headers = parameters
    .filter((parameter) => parameter.in === 'header')
    .map((parameter) => ({
      key: parameter.name,
      value: parameter.example ?? parameter.schema?.example ?? '',
      description: parameter.description,
      disabled: parameter.required !== true,
    }));
  const requestBody = makeRequestBody(operation, document);
  if (requestBody?.mode === 'raw') {
    headers.push({ key: 'Content-Type', value: 'application/json' });
  }

  const description = operation.description ?? operation.summary;
  const security =
    operation.security ??
    document.security ??
    [];
  const securitySchemeName = security[0] ? Object.keys(security[0])[0] : undefined;
  const securityScheme = securitySchemeName
    ? document.components?.securitySchemes?.[securitySchemeName]
    : undefined;
  const auth = makePostmanAuth(securityScheme);
  return {
    name: operation.summary ?? `${method.toUpperCase()} ${pathTemplate}`,
    ...(description ? { description } : {}),
    request: {
      method: method.toUpperCase(),
      header: headers,
      url: makeRequestUrl(baseUrl, pathTemplate, parameters),
      ...(auth ? { auth } : {}),
      ...(requestBody ? { body: requestBody } : {}),
    },
    response: [],
  };
}

function makePostmanAuth(scheme) {
  if (scheme?.type === 'http' && scheme.scheme === 'bearer') {
    return {
      type: 'bearer',
      bearer: [{ key: 'token', value: '{{accessToken}}', type: 'string' }],
    };
  }
  if (scheme?.type === 'http' && scheme.scheme === 'basic') {
    return {
      type: 'basic',
      basic: [
        { key: 'username', value: '{{username}}', type: 'string' },
        { key: 'password', value: '{{password}}', type: 'string' },
      ],
    };
  }
  if (
    scheme?.type === 'apiKey' &&
    (scheme.in === 'header' || scheme.in === 'query')
  ) {
    return {
      type: 'apikey',
      apikey: [
        { key: 'key', value: scheme.name, type: 'string' },
        { key: 'value', value: `{{${scheme.name}}}`, type: 'string' },
        { key: 'in', value: scheme.in, type: 'string' },
      ],
    };
  }
  return undefined;
}

function createPostmanCollection(document, baseUrl) {
  const folders = new Map();
  const operations = [];

  for (const [pathTemplate, pathItem] of Object.entries(document.paths ?? {})) {
    const pathParameters = pathItem.parameters ?? [];
    for (const [method, operation] of Object.entries(pathItem)) {
      if (!['get', 'put', 'post', 'delete', 'patch', 'head', 'options'].includes(method)) {
        continue;
      }

      const item = makePostmanItem(
        document,
        baseUrl,
        pathTemplate,
        method,
        operation,
        pathParameters,
      );
      const tag = operation.tags?.[0] ?? 'Cosmix API';
      if (!folders.has(tag)) {
        const folder = { name: tag, item: [] };
        folders.set(tag, folder);
        operations.push({ tag, folder });
      }
      folders.get(tag).item.push(item);
    }
  }

  return {
    info: {
      name: document.info?.title ?? 'OpenAPI collection',
      description: document.info?.description ?? '',
      schema: COLLECTION_SCHEMA,
    },
    item: operations.map(({ folder }) => folder),
    variable: [
      { key: 'baseUrl', value: baseUrl, type: 'string' },
      { key: 'accessToken', value: '', type: 'string' },
    ],
  };
}

function getOpenApiBaseUrl(document, fallbackUrl) {
  const server = document.servers?.[0];
  if (!server?.url) {
    return fallbackUrl;
  }

  return server.url.replace(/\{([^}]+)\}/g, (placeholder, name) => {
    const variable = server.variables?.[name];
    return variable?.default ?? variable?.enum?.[0] ?? placeholder;
  });
}

function normalizePath(pathname) {
  return pathname
    .replace(/:([A-Za-z0-9_]+)/g, '{$1}')
    .replace(/\/+$/, '') || '/';
}

function getRequestPath(request) {
  if (typeof request.url === 'string') {
    try {
      return normalizePath(new URL(request.url).pathname.replace(/:([A-Za-z0-9_]+)/g, '{$1}'));
    } catch {
      return normalizePath(request.url.split('?')[0]);
    }
  }

  if (request.url?.path) {
    return normalizePath(
      `/${request.url.path.join('/').replace(/:([A-Za-z0-9_]+)/g, '{$1}')}`,
    );
  }
  if (request.url?.raw) {
    try {
      return normalizePath(
        new URL(request.url.raw).pathname.replace(/:([A-Za-z0-9_]+)/g, '{$1}'),
      );
    } catch {
      return normalizePath(request.url.raw.split('?')[0]);
    }
  }
  return null;
}

function mergeCollection(existing, generated) {
  const endpoints = new Map();
  const generatedItems = [];
  existing.item ??= [];

  function collect(items) {
    for (const item of items ?? []) {
      if (item.request?.method && item.request?.url) {
        const key = `${item.request.method.toUpperCase()} ${getRequestPath(item.request)}`;
        endpoints.set(key, item);
      } else if (item.item) {
        collect(item.item);
      }
    }
  }

  collect(existing.item);
  for (const folder of generated.item) {
    for (const item of folder.item) {
      const key = `${item.request.method} ${getRequestPath(item.request)}`;
      const current = endpoints.get(key);
      if (current) {
        if (!item.request.auth && current.request.auth) {
          item.request.auth = current.request.auth;
        }
        current.name = item.name;
        current.request = item.request;
        current.description = item.description;
      } else {
        generatedItems.push({ folderName: folder.name, item });
      }
    }
  }

  for (const { folderName, item } of generatedItems) {
    let folder = existing.item.find(
      (entry) => entry.name === folderName && Array.isArray(entry.item),
    );
    if (!folder) {
      folder = { name: folderName, item: [] };
      existing.item.push(folder);
    }
    folder.item.push(item);
  }

  existing.variable ??= [];
  for (const variable of generated.variable ?? []) {
    if (!existing.variable.some((current) => current.key === variable.key)) {
      existing.variable.push(variable);
    }
  }

  existing.info = { ...existing.info, schema: COLLECTION_SCHEMA };
  return existing;
}

async function postmanRequest(apiKey, endpoint, options = {}) {
  const response = await fetch(`${POSTMAN_API_URL}${endpoint}`, {
    ...options,
    headers: {
      'X-Api-Key': apiKey,
      Accept: 'application/json',
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
      ...options.headers,
    },
  });
  const body = await response.text();
  if (!response.ok) {
    throw new Error(
      `Postman API ${options.method ?? 'GET'} ${endpoint} failed (${response.status}): ${body}`,
    );
  }
  return body ? JSON.parse(body) : {};
}

async function listCollections(apiKey) {
  const collections = [];
  for (let offset = 0; ; offset += 100) {
    const result = await postmanRequest(
      apiKey,
      `/collections?limit=100&offset=${offset}`,
    );
    const page = result.collections ?? [];
    collections.push(...page);
    if (page.length < 100) {
      return collections;
    }
  }
}

async function syncPostman() {
  loadEnvFile(path.join(ROOT_DIR, 'services', 'gateway-service', '.env.local'));

  const apiKey = process.env.POSTMAN_API_KEY;
  const collectionName = process.env.POSTMAN_COLLECTION_NAME;
  const swaggerUrl =
    process.env.SWAGGER_URL ?? 'http://localhost:5000/api/docs-json';
  const baseUrlOverride = process.env.POSTMAN_BASE_URL;
  const workspaceId = process.env.POSTMAN_WORKSPACE_ID;

  if (!apiKey || !collectionName) {
    throw new Error(
      'Set POSTMAN_API_KEY and POSTMAN_COLLECTION_NAME in services/gateway-service/.env.local.',
    );
  }

  const swaggerResponse = await fetch(swaggerUrl);
  if (!swaggerResponse.ok) {
    throw new Error(
      `Could not fetch OpenAPI document from ${swaggerUrl} (${swaggerResponse.status}).`,
    );
  }
  const document = await swaggerResponse.json();
  if (!document.openapi || !document.paths) {
    throw new Error('The Swagger URL did not return a valid OpenAPI document.');
  }

  const baseUrl =
    baseUrlOverride ??
    getOpenApiBaseUrl(document, new URL(swaggerUrl).origin);
  const generated = createPostmanCollection(document, baseUrl);
  const collections = await listCollections(apiKey);
  const matches = collections.filter((collection) => collection.name === collectionName);

  if (matches.length > 1) {
    throw new Error(
      `Found multiple Postman collections named "${collectionName}". Rename or remove duplicates before syncing.`,
    );
  }

  if (matches.length === 0) {
    const query = workspaceId ? `?workspaceId=${encodeURIComponent(workspaceId)}` : '';
    const result = await postmanRequest(apiKey, `/collections${query}`, {
      method: 'POST',
      body: JSON.stringify({ collection: { ...generated, info: { ...generated.info, name: collectionName } } }),
    });
    console.log(`Created Postman collection "${collectionName}" (${result.collection?.uid ?? 'created'}).`);
    return;
  }

  const uid = matches[0].uid;
  if (!uid) {
    throw new Error(`Postman did not return a collection UID for "${collectionName}".`);
  }
  const result = await postmanRequest(apiKey, `/collections/${encodeURIComponent(uid)}`);
  const existing = result.collection;
  if (!existing) {
    throw new Error(`Postman returned no collection data for "${collectionName}".`);
  }

  const previousCount = countRequests(existing.item);
  const merged = mergeCollection(existing, generated);
  await postmanRequest(apiKey, `/collections/${encodeURIComponent(uid)}`, {
    method: 'PUT',
    body: JSON.stringify({ collection: merged }),
  });
  console.log(
    `Synced "${collectionName}": ${countRequests(generated.item)} Swagger endpoints processed; unrelated requests were preserved.`,
  );
  console.log(`Existing collection contained ${previousCount} requests before sync.`);
}

function countRequests(items) {
  return (items ?? []).reduce(
    (count, item) =>
      count + (item.request ? 1 : 0) + (item.item ? countRequests(item.item) : 0),
    0,
  );
}

if (require.main === module) {
  syncPostman().catch((error) => {
    console.error(`Postman sync failed: ${error.message}`);
    process.exitCode = 1;
  });
}

module.exports = {
  createPostmanCollection,
  getOpenApiBaseUrl,
  makePostmanAuth,
  mergeCollection,
  normalizePath,
  getRequestPath,
};
