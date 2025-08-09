// src/lib/api-error.ts
export type FieldErrors = Record<string, string>;

export type ParsedApiError = {
  fieldErrors: FieldErrors;
  formError: string | null;
};

// Narrowing helpers
function isResponse(x: unknown): x is Response {
  return typeof x === "object" && x !== null && "ok" in x && "status" in x;
}

export async function parseApiError(
  input: Response | unknown
): Promise<ParsedApiError> {
  // Network/timeout/offline cases (fetch threw before giving us a Response)
  if (!isResponse(input)) {
    // Timeout via AbortController usually sets name to "AbortError"

    const e = input as Error;

    if (typeof navigator !== "undefined" && !navigator.onLine) {
      return {
        fieldErrors: {},
        formError:
          "You appear to be offline. Check your connection and try again.",
      };
    }
    if (e?.name === "AbortError") {
      return {
        fieldErrors: {},
        formError: "The request timed out. Please try again.",
      };
    }
    return { fieldErrors: {}, formError: "Network error. Please try again." };
  }

  const res = input;

  const fallback: ParsedApiError = {
    fieldErrors: {},
    formError: `Unexpected error (${res.status})`,
  };

  const raw = await res.text();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let data: any;
  try {
    data = raw ? JSON.parse(raw) : {};
  } catch {
    return { fieldErrors: {}, formError: raw || fallback.formError };
  }

  // 422 from validator
  if (data.FieldErrors || data.Errors) {
    return {
      fieldErrors: (data.FieldErrors ?? {}) as FieldErrors,
      formError:
        Array.isArray(data.Errors) && data.Errors.length
          ? data.Errors[0]
          : null,
    };
  }

  // 400/404/405/500/503 message shape from backend
  if (typeof data.Error === "string") {
    return { fieldErrors: {}, formError: data.Error };
  }

  return fallback;
}
