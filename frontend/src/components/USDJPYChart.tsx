import React from 'react';
import { USDJPYRate } from '../types/trade';

interface USDJPYChartProps {
  points: USDJPYRate[];
  timeframe: 'daily' | 'weekly';
}

const WIDTH = 980;
const HEIGHT = 320;
const LEFT_PADDING = 72;
const RIGHT_PADDING = 40;
const TOP_PADDING = 28;
const BOTTOM_PADDING = 40;

const pickNiceStep = (rawStep: number): number => {
  const exponent = Math.floor(Math.log10(rawStep));
  const fraction = rawStep / Math.pow(10, exponent);

  if (fraction <= 1) return Math.pow(10, exponent);
  if (fraction <= 2) return 2 * Math.pow(10, exponent);
  if (fraction <= 5) return 5 * Math.pow(10, exponent);
  return 10 * Math.pow(10, exponent);
};

const calcTickValues = (min: number, max: number, targetTickCount = 7): number[] => {
  const range = Math.max(max - min, 0.0001);
  const rawStep = range / Math.max(targetTickCount - 1, 1);
  const step = pickNiceStep(rawStep);

  const start = Math.floor(min / step) * step;
  const end = Math.ceil(max / step) * step;

  const ticks: number[] = [];
  for (let value = start; value <= end + step * 0.5; value += step) {
    ticks.push(Number(value.toFixed(6)));
  }
  return ticks;
};

const formatRate = (value: number, step: number): string => {
  if (step >= 1) return value.toFixed(0);
  if (step >= 0.1) return value.toFixed(1);
  if (step >= 0.01) return value.toFixed(2);
  if (step >= 0.001) return value.toFixed(3);
  return value.toFixed(4);
};

const USDJPYChart: React.FC<USDJPYChartProps> = ({ points, timeframe }) => {
  const [hoveredIndex, setHoveredIndex] = React.useState<number | null>(null);

  if (points.length === 0) {
    return <div style={styles.message}>表示できるレートデータがありません</div>;
  }

  const closePrices = points.map((p) => p.close);
  const min = Math.min(...closePrices);
  const max = Math.max(...closePrices);
  const tickValues = calcTickValues(min, max);
  const tickStep = tickValues.length > 1 ? tickValues[1] - tickValues[0] : 0.001;
  const top = tickValues[tickValues.length - 1];
  const bottom = tickValues[0];
  const chartWidth = WIDTH - LEFT_PADDING - RIGHT_PADDING;
  const chartHeight = HEIGHT - TOP_PADDING - BOTTOM_PADDING;

  const xForIndex = (index: number) => {
    if (points.length === 1) {
      return WIDTH / 2;
    }
    return LEFT_PADDING + (chartWidth * index) / (points.length - 1);
  };

  const yForValue = (value: number) => {
    return TOP_PADDING + ((top - value) / (top - bottom)) * chartHeight;
  };

  const polylinePoints = points
    .map((point, index) => `${xForIndex(index)},${yForValue(point.close)}`)
    .join(' ');

  const handleMouseMove = (event: React.MouseEvent<SVGSVGElement>) => {
    const rect = event.currentTarget.getBoundingClientRect();
    if (rect.width === 0) return;

    const mouseX = ((event.clientX - rect.left) / rect.width) * WIDTH;
    const clampedX = Math.min(Math.max(mouseX, LEFT_PADDING), WIDTH - RIGHT_PADDING);

    let nearest = 0;
    let nearestDistance = Math.abs(xForIndex(0) - clampedX);
    for (let i = 1; i < points.length; i += 1) {
      const distance = Math.abs(xForIndex(i) - clampedX);
      if (distance < nearestDistance) {
        nearest = i;
        nearestDistance = distance;
      }
    }

    setHoveredIndex(nearest);
  };

  const handleMouseLeave = () => {
    setHoveredIndex(null);
  };

  const firstLabel = points[0].label || points[0].date;
  const lastLabel = points[points.length - 1].label || points[points.length - 1].date;

  const hoveredPoint = hoveredIndex !== null ? points[hoveredIndex] : null;
  const hoveredX = hoveredIndex !== null ? xForIndex(hoveredIndex) : 0;
  const hoveredY = hoveredPoint ? yForValue(hoveredPoint.close) : 0;

  const tooltipWidth = 170;
  const tooltipHeight = 112;
  const tooltipX = hoveredPoint
    ? Math.min(Math.max(hoveredX + 12, LEFT_PADDING + 8), WIDTH - RIGHT_PADDING - tooltipWidth)
    : 0;
  const tooltipY = hoveredPoint
    ? Math.min(Math.max(hoveredY - tooltipHeight - 12, TOP_PADDING + 8), HEIGHT - BOTTOM_PADDING - tooltipHeight)
    : 0;

  return (
    <div style={styles.wrapper}>
      <svg
        viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
        width="100%"
        height="auto"
        role="img"
        aria-label="USDJPY close chart"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
      >
        <rect x={0} y={0} width={WIDTH} height={HEIGHT} fill="#ffffff" rx={10} />

        <line x1={LEFT_PADDING} y1={TOP_PADDING} x2={LEFT_PADDING} y2={HEIGHT - BOTTOM_PADDING} stroke="#cfd6de" strokeWidth={1} />
        <line
          x1={LEFT_PADDING}
          y1={HEIGHT - BOTTOM_PADDING}
          x2={WIDTH - RIGHT_PADDING}
          y2={HEIGHT - BOTTOM_PADDING}
          stroke="#cfd6de"
          strokeWidth={1}
        />

        {tickValues.map((tick, idx) => {
          const y = yForValue(tick);
          const isBoundary = idx === 0 || idx === tickValues.length - 1;
          return (
            <g key={`tick-${tick}`}>
              <line
                x1={LEFT_PADDING}
                y1={y}
                x2={WIDTH - RIGHT_PADDING}
                y2={y}
                stroke={isBoundary ? '#d6dde6' : '#eef2f6'}
                strokeWidth={1}
              />
              <text x={LEFT_PADDING - 8} y={y + 4} fontSize="12" fill="#59636e" textAnchor="end">
                {formatRate(tick, tickStep)}
              </text>
            </g>
          );
        })}

        <polyline fill="none" stroke="#0d6efd" strokeWidth={3} points={polylinePoints} />

        {points.map((point, index) => (
          <circle
            key={`${point.date}-${index}`}
            cx={xForIndex(index)}
            cy={yForValue(point.close)}
            r={3}
            fill="#0d6efd"
          />
        ))}

        {hoveredPoint && (
          <g>
            <line
              x1={hoveredX}
              y1={TOP_PADDING}
              x2={hoveredX}
              y2={HEIGHT - BOTTOM_PADDING}
              stroke="#8aa8d8"
              strokeWidth={1}
              strokeDasharray="4 4"
            />
            <circle cx={hoveredX} cy={hoveredY} r={5} fill="#0d6efd" stroke="#ffffff" strokeWidth={2} />

            <rect x={tooltipX} y={tooltipY} width={tooltipWidth} height={tooltipHeight} rx={8} fill="#ffffff" stroke="#cfd6de" />
            <text x={tooltipX + 10} y={tooltipY + 18} fontSize="12" fill="#1f2937">{hoveredPoint.label || hoveredPoint.date}</text>
            <text x={tooltipX + 10} y={tooltipY + 36} fontSize="12" fill="#374151">{`Bid: ${hoveredPoint.bid.toFixed(3)}`}</text>
            <text x={tooltipX + 10} y={tooltipY + 52} fontSize="12" fill="#374151">{`Ask: ${hoveredPoint.ask.toFixed(3)}`}</text>
            <text x={tooltipX + 10} y={tooltipY + 68} fontSize="12" fill="#374151">{`High: ${hoveredPoint.high.toFixed(3)}`}</text>
            <text x={tooltipX + 10} y={tooltipY + 84} fontSize="12" fill="#374151">{`Low: ${hoveredPoint.low.toFixed(3)}`}</text>
            <text x={tooltipX + 10} y={tooltipY + 100} fontSize="12" fill="#111827">{`Close: ${hoveredPoint.close.toFixed(3)}`}</text>
          </g>
        )}

        <text x={LEFT_PADDING} y={16} fontSize="13" fill="#2b2f33">{`High ${max.toFixed(3)}`}</text>
        <text x={LEFT_PADDING + 130} y={16} fontSize="13" fill="#2b2f33">{`Low ${min.toFixed(3)}`}</text>
        <text x={LEFT_PADDING} y={HEIGHT - 10} fontSize="12" fill="#59636e">{firstLabel}</text>
        <text x={WIDTH - RIGHT_PADDING} y={HEIGHT - 10} fontSize="12" fill="#59636e" textAnchor="end">{lastLabel}</text>
        <text x={WIDTH / 2} y={16} fontSize="13" fill="#59636e" textAnchor="middle">{timeframe === 'weekly' ? '週足' : '日足'}</text>
      </svg>
    </div>
  );
};

const styles = {
  wrapper: {
    backgroundColor: '#ffffff',
    borderRadius: '10px',
    border: '1px solid #e2e8f0',
    padding: '12px',
    boxShadow: '0 2px 6px rgba(0, 0, 0, 0.06)',
  },
  message: {
    padding: '20px',
    textAlign: 'center' as const,
    color: '#6c757d',
    backgroundColor: '#ffffff',
    borderRadius: '8px',
    border: '1px solid #e2e8f0',
  },
};

export default USDJPYChart;
