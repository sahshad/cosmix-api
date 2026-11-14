import { Injectable, LoggerService } from '@nestjs/common';

type LogLevel = 'debug' | 'info' | 'warn' | 'error';

type LogFields = Record<string, unknown>;

const SERVICE_NAME = 'gateway-service';

const RESET = '\x1b[0m';
const BLUE = '\x1b[34m';
const YELLOW = '\x1b[33m';
const RED = '\x1b[31m';
const PURPLE = '\x1b[35m';
const CYAN = '\x1b[36m';
const GREEN = '\x1b[32m';
const GRAY = '\x1b[90m';

const LEVEL_ORDER: Record<LogLevel, number> = {
  debug: 10,
  info: 20,
  warn: 30,
  error: 40,
};

const LEVEL_DISPLAY: Record<LogLevel, { color: string; label: string }> = {
  debug: { color: GRAY, label: 'DEBUG' },
  info: { color: BLUE, label: 'INFO' },
  warn: { color: YELLOW, label: 'WARN' },
  error: { color: RED, label: 'ERROR' },
};

const FIELD_ORDER = [
  'method',
  'path',
  'service',
  'status',
  'request_id',
  'user_id',
  'grpc_code',
  'duration_ms',
  'error',
  'context',
  'stack',
];

const FIELD_LABELS: Record<string, string> = {
  method: 'Method',
  path: 'Path',
  service: 'Service',
  status: 'Status',
  request_id: 'Request ID',
  user_id: 'User ID',
  grpc_code: 'gRPC Code',
  duration_ms: 'Duration',
  error: 'Error',
  context: 'Context',
  stack: 'Stack',
};

const FIELD_COLORS: Record<string, string> = {
  method: CYAN,
  path: CYAN,
  service: PURPLE,
  request_id: GRAY,
  user_id: GRAY,
  grpc_code: YELLOW,
  duration_ms: GREEN,
  error: RED,
  context: GRAY,
  stack: GRAY,
};

const NOISY_CONTEXTS = new Set([
  'InstanceLoader',
  'RouterExplorer',
  'RoutesResolver',
]);

function isProduction(): boolean {
  return (
    process.env.NODE_ENV === 'production' ||
    process.env.APP_ENV === 'production'
  );
}

function resolveLogLevel(): LogLevel {
  const configured = process.env.LOG_LEVEL?.toLowerCase();

  if (configured && configured in LEVEL_ORDER) {
    return configured as LogLevel;
  }

  return isProduction() ? 'info' : 'debug';
}

function isPlainObject(value: unknown): value is LogFields {
  return (
    typeof value === 'object' &&
    value !== null &&
    !Array.isArray(value) &&
    !(value instanceof Error)
  );
}

function messageFrom(value: unknown): string {
  if (typeof value === 'string') {
    return value;
  }

  if (value instanceof Error) {
    return value.message;
  }

  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

function pad(value: number, length = 2): string {
  return String(value).padStart(length, '0');
}

function timestamp(): string {
  const now = new Date();

  return `${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}.${pad(now.getMilliseconds(), 3)}`;
}

function normalize(
  message: unknown,
  params: unknown[],
): { text: string; fields: LogFields } {
  const fields: LogFields = {};

  for (const param of params) {
    if (param === undefined || param === null) {
      continue;
    }

    if (typeof param === 'string') {
      if (param.includes('\n')) {
        fields.stack = param;
      } else {
        fields.context = param;
      }
      continue;
    }

    if (param instanceof Error) {
      fields.error = param.message;

      if (param.stack) {
        fields.stack = param.stack;
      }
      continue;
    }

    if (isPlainObject(param)) {
      Object.assign(fields, param);
      continue;
    }

    fields.detail = messageFrom(param);
  }

  return { text: messageFrom(message), fields };
}

function serializeFields(fields: LogFields): LogFields {
  const serialized: LogFields = {};

  for (const [key, value] of Object.entries(fields)) {
    if (value === undefined) {
      continue;
    }

    serialized[key] = value instanceof Error ? value.message : value;
  }

  return serialized;
}

function orderedEntries(fields: LogFields): [string, unknown][] {
  const entries: [string, unknown][] = [];

  for (const key of FIELD_ORDER) {
    if (fields[key] !== undefined) {
      entries.push([key, fields[key]]);
    }
  }

  for (const key of Object.keys(fields).sort()) {
    if (!FIELD_ORDER.includes(key) && fields[key] !== undefined) {
      entries.push([key, fields[key]]);
    }
  }

  return entries;
}

function fieldLabel(key: string): string {
  if (FIELD_LABELS[key]) {
    return FIELD_LABELS[key];
  }

  return key
    .split('_')
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ');
}

function fieldColor(key: string, value: unknown): string | undefined {
  if (key === 'status' && typeof value === 'number') {
    if (value >= 500) {
      return RED;
    }

    if (value >= 400) {
      return YELLOW;
    }

    if (value >= 300) {
      return CYAN;
    }

    return GREEN;
  }

  return FIELD_COLORS[key];
}

function formatValue(key: string, value: unknown): string {
  if (key === 'duration_ms' && typeof value === 'number') {
    return `${value}ms`;
  }

  if (typeof value === 'string') {
    return value;
  }

  if (
    typeof value === 'number' ||
    typeof value === 'boolean' ||
    value === null
  ) {
    return String(value);
  }

  if (value instanceof Error) {
    return value.message;
  }

  try {
    return JSON.stringify(value);
  } catch {
    return Object.prototype.toString.call(value) as string;
  }
}

@Injectable()
export class AppLogger implements LoggerService {
  log(message: unknown, ...optionalParams: unknown[]): void {
    this.write('info', message, optionalParams);
  }

  info(message: unknown, ...optionalParams: unknown[]): void {
    this.write('info', message, optionalParams);
  }

  error(message: unknown, ...optionalParams: unknown[]): void {
    this.write('error', message, optionalParams);
  }

  warn(message: unknown, ...optionalParams: unknown[]): void {
    this.write('warn', message, optionalParams);
  }

  debug(message: unknown, ...optionalParams: unknown[]): void {
    this.write('debug', message, optionalParams);
  }

  verbose(message: unknown, ...optionalParams: unknown[]): void {
    this.write('debug', message, optionalParams);
  }

  fatal(message: unknown, ...optionalParams: unknown[]): void {
    this.write('error', message, optionalParams);
  }

  private write(
    level: LogLevel,
    message: unknown,
    optionalParams: unknown[],
  ): void {
    const { text, fields } = normalize(message, optionalParams);
    const context =
      typeof fields.context === 'string' ? fields.context : undefined;

    if (
      context &&
      NOISY_CONTEXTS.has(context) &&
      LEVEL_ORDER[level] < LEVEL_ORDER.warn
    ) {
      return;
    }

    if (LEVEL_ORDER[level] < LEVEL_ORDER[resolveLogLevel()]) {
      return;
    }

    if (isProduction()) {
      process.stdout.write(
        `${JSON.stringify({
          level,
          ts: Date.now() / 1000,
          logger: SERVICE_NAME,
          msg: text,
          ...serializeFields(fields),
        })}\n`,
      );
      return;
    }

    process.stdout.write(this.pretty(level, text, fields) + '\n');
  }

  private pretty(
    level: LogLevel,
    message: string,
    rawFields: LogFields,
  ): string {
    const fields = serializeFields(rawFields);
    const display = LEVEL_DISPLAY[level];
    const rows = orderedEntries(fields).map(([key, value]) => ({
      label: fieldLabel(key),
      value: formatValue(key, value),
      color: fieldColor(key, value),
    }));
    const labelWidth =
      rows.reduce((width, row) => Math.max(width, row.label.length), 0) + 2;

    const lines = [
      `${timestamp()}  ${display.color}${display.label.padEnd(5)}${RESET} [${PURPLE}${SERVICE_NAME}${RESET}] ${message}`,
    ];

    for (const row of rows) {
      const value = row.color ? `${row.color}${row.value}${RESET}` : row.value;

      lines.push(
        `    ${row.label}${' '.repeat(labelWidth - row.label.length)}: ${value}`,
      );
    }

    return lines.join('\n') + '\n';
  }
}

export const logger = new AppLogger();
