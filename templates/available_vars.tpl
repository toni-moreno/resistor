WINDOW_NODE='\n    \|window()\n        .period(INTERVAL_CHECK)\n        .every(EVERY)\n        .align()'
AVAILABLE_FUNCTIONS=("last" "max" "mean" "median" "min" "spread" "stddev" "sum" "movingAverage" "percentile")
TEMPLATE_TYPE=("GAUGE" "COUNTER")
DERIV_NODE="\n    \|derivative(\'value\')\n        .unit(UNIT_DERIV)\n        .nonNegative()"
AVAILABLE_TH_DIRECTIONS=(">" "<")
AVAILABLE_TXX_AC_DIRECTIONS=(">" "<")
AVAILABLE_TXX_DC_DIRECTIONS=("<" ">")
AVAILABLE_TXP_DC_DIRECTIONS0=">"
AVAILABLE_TXN_DC_DIRECTIONS0="<"
TREND_TYPE=("" '/float("past.value") * 100.0')