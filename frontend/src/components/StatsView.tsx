import fetchStats from "@/api/statsAPI";
import { useQuery } from "@tanstack/react-query";
import CalendarHeatmap from "react-calendar-heatmap";
import "react-calendar-heatmap/dist/styles.css";
import { Tooltip } from "react-tooltip";
import {
  ChartContainer,
  type ChartConfig,
  ChartTooltip,
  ChartTooltipContent,
} from "./ui/chart";
import { BarChart, Bar, CartesianGrid, XAxis } from "recharts";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectLabel,
  SelectItem,
} from "./ui/select";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardAction,
  CardContent,
} from "./ui/card";
import { ToggleGroup, ToggleGroupItem } from "./ui/toggle-group";
import { useState } from "react";

interface HeatMapCell {
  date: Date;
  count: number;
}

enum Duration {
  Week = "6",
  Month = "29",
  Year = "364",
}

export default function StatsView() {
  const [forecastTimeRange, setForecastTimeRange] = useState<string>(
    Duration.Week,
  );

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["stats", forecastTimeRange],
    queryFn: () => fetchStats(forecastTimeRange),
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

  // (AJama) Forecast Chart
  const chartConfig = {
    Due: {
      label: "Due Cards",
      color: "#2563eb",
    },
  } satisfies ChartConfig;
  console.log(data.stats);

  return (
    <div className="@container/main flex flex-1 flex-col gap-2">
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
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
        <Card className="@container/card mx-5 px-4 lg:px-6">
          <CardHeader>
            <CardTitle>Forecast</CardTitle>
            <CardDescription>
              <span className="hidden @[540px]/card:block">
                Cards due in the future.
              </span>
              <span className="@[540px]/card:hidden">
                Cards due in the future.
              </span>
            </CardDescription>
            <CardAction>
              <ToggleGroup
                type="single"
                value={forecastTimeRange}
                onValueChange={setForecastTimeRange}
                variant="outline"
                className="hidden *:data-[slot=toggle-group-item]:px-4! @[767px]/card:flex"
              >
                <ToggleGroupItem value={Duration.Week}>
                  Next Week
                </ToggleGroupItem>
                <ToggleGroupItem value={Duration.Month}>
                  Next Month
                </ToggleGroupItem>
                <ToggleGroupItem value={Duration.Year}>
                  Next Year
                </ToggleGroupItem>
              </ToggleGroup>
              <Select
                value={forecastTimeRange}
                onValueChange={setForecastTimeRange}
              >
                <SelectTrigger
                  className="flex w-40 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[767px]/card:hidden"
                  aria-label="Select a value"
                >
                  <SelectValue placeholder="Last 3 months" />
                </SelectTrigger>
                <SelectContent className="rounded-xl">
                  <SelectItem value={Duration.Week} className="rounded-lg">
                    Next Week
                  </SelectItem>
                  <SelectItem value={Duration.Month} className="rounded-lg">
                    Next Month
                  </SelectItem>
                  <SelectItem value={Duration.Year} className="rounded-lg">
                    Next Year
                  </SelectItem>
                </SelectContent>
              </Select>
            </CardAction>
          </CardHeader>
          <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
            <ChartContainer
              config={chartConfig}
              className="aspect-auto h-[250px] w-full"
            >
              <BarChart accessibilityLayer data={data.stats.forecasts}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="due_date"
                  tickLine={false}
                  tickMargin={10}
                  axisLine={false}
                  tickFormatter={(value) => value.slice(0, 10)}
                />
                <ChartTooltip
                  content={
                    <ChartTooltipContent
                      labelFormatter={(value) => {
                        return new Date(value).toLocaleDateString("en-US", {
                          month: "short",
                          day: "numeric",
                        });
                      }}
                      nameKey="Due"
                    />
                  }
                />
                <Bar dataKey="due_count" fill="#2563eb" radius={4} />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
