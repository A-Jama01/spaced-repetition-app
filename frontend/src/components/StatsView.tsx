import fetchStats from "@/api/statsAPI";
import { useQuery } from "@tanstack/react-query";
import CalendarHeatmap from "react-calendar-heatmap";
import "react-calendar-heatmap/dist/styles.css";
import { Tooltip } from "react-tooltip";

interface HeatMapCell {
  date: Date;
  count: number;
}

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

  // (AJama) Get values to populate heatmap
  const YEAR_SIZE = 366;
  const currentDate = new Date();
  const startDate = new Date(currentDate);
  startDate.setFullYear(startDate.getFullYear() - 1);

  const heatmapMap = new Map();
  for (let i = 0; i < YEAR_SIZE; i++) {
    const nextDate = new Date(currentDate);
    nextDate.setDate(nextDate.getDate() - i);
    heatmapMap.set(nextDate, 0);
  }

  for (let i = 0; i < data.stats.heatmap.length; i++) {
    const heatmapCellDate = new Date(data.stats.heatmap[i].review_date);
    heatmapMap.set(heatmapCellDate, data.stats.heatmap[i].reviews);
  }

  const heatmapValues: HeatMapCell[] = [];
  heatmapMap.forEach((value: number, key: Date) => {
    heatmapValues.push({
      date: key,
      count: value,
    });
  });

  const maxHeatmapValue = Math.max(...heatmapValues.map((cell) => cell.count));

  function classForValue(
    value:
      | CalendarHeatmap.ReactCalendarHeatmapValue<CalendarHeatmap.ReactCalendarHeatmapDate>
      | undefined,
  ): string {
    if (!value || value.count === 0) {
      return "color-scale-0";
    }

    const level = Math.ceil((value.count / maxHeatmapValue) * 4);
    return `color-scale-${level}`;
  }

  function getTooltipDataAttrs(
    value:
      | CalendarHeatmap.ReactCalendarHeatmapValue<CalendarHeatmap.ReactCalendarHeatmapDate>
      | undefined,
  ): any {
    if (!value) return null;

    const date = new Date(value.date).toLocaleDateString();

    return {
      "data-tooltip-id": "heatmap-tooltip",
      "data-tooltip-content": `${date}: ${value.count} reviews`,
    };
  }

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
      <div className="flex flex-row">
        <CalendarHeatmap
          startDate={startDate}
          endDate={currentDate}
          values={heatmapValues}
          showWeekdayLabels={true}
          classForValue={classForValue}
          tooltipDataAttrs={getTooltipDataAttrs}
        />
        <Tooltip id="heatmap-tooltip" />
      </div>
    </div>
  );
}
