import { NextRequest, NextResponse } from "next/server";

const BASE = process.env.API_BASE_URL!;

export async function GET(
  _: NextRequest,
  { params }: { params: { id: string } }
) {
  const res = await fetch(`${BASE}/books/${params.id}`, { cache: "no-store" });
  const body = await res.text();
  return new NextResponse(body, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}

export async function PUT(
  req: NextRequest,
  { params }: { params: { id: string } }
) {
  const json = await req.text();
  const res = await fetch(`${BASE}/books/${params.id}`, {
    method: "PUT",
    body: json,
    headers: { "Content-Type": "application/json" },
  });
  const body = await res.text();
  return new NextResponse(body, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}

export async function DELETE(
  _: NextRequest,
  { params }: { params: { id: string } }
) {
  const res = await fetch(`${BASE}/books/${params.id}`, { method: "DELETE" });
  const body = await res.text();
  return new NextResponse(body || null, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}
