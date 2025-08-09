export async function apiGet(path: string, init?: RequestInit) {
  const base = process.env.API_BASE_URL!;
  const res = await fetch(`${base}${path}`, {
    cache: "no-store",
    ...init,
  });
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}
