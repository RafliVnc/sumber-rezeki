type Method = "GET" | "POST" | "PUT" | "DELETE";

interface ApiOptions<TBody = unknown> {
  url: string;
  method?: Method;
  body?: TBody;
  params?: Record<
    string,
    string | number | boolean | Array<string | number> | undefined | null
  >;
}

const ApiUrl = process.env.NEXT_PUBLIC_API_URL;

export const api = async <T, TBody = unknown>({
  url,
  method = "GET",
  body,
  params,
}: ApiOptions<TBody>): Promise<T> => {
  const token = localStorage.getItem("token");

  // Convert params ke query string dengan support array
  let queryString = "";
  if (params) {
    const searchParams = new URLSearchParams();

    Object.entries(params).forEach(([key, value]) => {
      if (value === undefined || value === null) return;

      // Handle array params
      if (Array.isArray(value)) {
        value.forEach((v) => searchParams.append(`${key}[]`, String(v)));
      } else {
        searchParams.append(key, String(value));
      }
    });

    const queryStr = searchParams.toString();
    if (queryStr) {
      queryString = "?" + queryStr;
    }
  }

  const res = await fetch(`${ApiUrl}/${url}${queryString}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `${token}` } : {}),
    },
    ...(body && method !== "GET" ? { body: JSON.stringify(body) } : {}),
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({}));
    throw new Error(error.message || `API Error: ${res.status}`);
  }

  return res.json() as Promise<T>;
};
