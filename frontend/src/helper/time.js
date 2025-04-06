export function parseDurationToSeconds(duration) {
    if (typeof duration !== "string") return "0.0000";
  
    const match = duration.match(/^([\d.]+)(ms|s|µs|ns)$/);
    if (!match) return "0.0000";
  
    const value = parseFloat(match[1]);
    const unit = match[2];
  
    switch (unit) {
      case "s":
        return value.toFixed(4);
      case "ms":
        return (value / 1000).toFixed(4);
      case "µs":
        return (value / 1e6).toFixed(4);
      case "ns":
        return (value / 1e9).toFixed(4);
      default:
        return "0.0000";
    }
  }
  