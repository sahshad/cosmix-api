import { handleGrpcError } from './grpc-error.util';

export function grpcUnaryCall<
  TResponse,
>(
  fn: (
    callback: (
      err: any,
      response: TResponse,
    ) => void,
  ) => void,
): Promise<TResponse> {
  return new Promise(
    (resolve, reject) => {
      fn((err, response) => {
        if (err) {
          return reject(
            handleGrpcError(err),
          );
        }

        resolve(response);
      });
    },
  );
}