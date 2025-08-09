import Skeleton from "@/components/Skeleton";

export default function LoadingBooks() {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <Skeleton className="h-7 w-40" />
        <Skeleton className="h-9 w-28" />
      </div>
      <div className="grid gap-3">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="rounded-lg border bg-white p-4">
            <Skeleton className="h-5 w-44" />
            <div className="mt-2 flex gap-2">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-12" />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
