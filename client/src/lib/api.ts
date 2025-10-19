type Method = "GET" | "POST" | "PUT" | "DELETE";

interface ApiOptions {
  url: string;
  method?: Method;
  body?: Record<string, any>;
  params?: Record<string, any>;
}

const ApiUrl = process.env.NEXT_PUBLIC_API_URL;

export const api = async <T>({
  url,
  method = "GET",
  body,
  params,
}: ApiOptions): Promise<T> => {
  const token = localStorage.getItem("token");

  // Convert params ke query string
  const queryString = params
    ? "?" +
      new URLSearchParams(
        Object.entries(params).reduce((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {} as Record<string, string>)
      ).toString()
    : "";

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
