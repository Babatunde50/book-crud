import { NextRequest, NextResponse } from "next/server";

const BASE = process.env.API_BASE_URL!;

export async function GET() {
  const res = await fetch(`${BASE}/books`, { cache: "no-store" });
  const body = await res.text();
  return new NextResponse(body, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}

export async function POST(req: NextRequest) {
  const json = await req.text();
  const res = await fetch(`${BASE}/books`, {
    method: "POST",
    body: json,
    headers: { "Content-Type": "application/json" },
  });
  const body = await res.text();
  return new NextResponse(body, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}
