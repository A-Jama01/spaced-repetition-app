import fetchStats from "@/api/statsAPI";
import { useQuery } from "@tanstack/react-query";
import CalendarHeatmap from "react-calendar-heatmap";
import "react-calendar-heatmap/dist/styles.css";

export default function StatsView() {
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["stats"],
    queryFn: fetchStats,
  });
  if (isPending) {
    return <span>Loading...</span>;
  }

  if (isError) {
    return <span>Error: {error.message}</span>;
  }

  console.log(data.stats);

  return (
    <div className="flex flex-1 flex-col">
      <h1 className="scroll-m-20 text-2xl text-center font-semibold tracking-tight">
        Today
      </h1>
      <div className="flex flex-row justify-evenly">
        <div className="flex flex-col items-center">
          <h2>Reviews</h2>
          <div>{data.stats.review_count}</div>
        </div>
        <div className="flex flex-col items-center">
          <h2>Retention Rate</h2>
          <div>{data.stats.retention}%</div>
        </div>
      </div>
      <h1 className="scroll-m-20 text-2xl text-center font-semibold tracking-tight">
        Heatmap
      </h1>
      <div className="flex flex-row justify-evenly"></div>
    </div>
  );
}
