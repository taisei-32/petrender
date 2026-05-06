#!/bin/bash

DATA="$1"
RENDER_DIR="render_imgs/$DATA"
ANALYZE_ZBAR_DIR="analyze/$DATA/zbar"
ANALYZE_ZXING_DIR="analyze/$DATA/zxing"
ZBARIMG_DIR="/opt/zbar/zbarimg/zbarimg"
ZXING_DIR="/opt/zxing-cpp/build/example/ZXingReader"
LOG_RAW_ZBAR="analyze/result_log/raw/zbar/${DATA}_zbar_log.txt"
LOG_RAW_ZXING="analyze/result_log/raw/zxing/${DATA}_zxing_log.txt"
LOG_REG_ZBAR="analyze/result_log/reg/zbar/${DATA}_zbar_log.txt"
LOG_REG_ZXING="analyze/result_log/reg/zxing/${DATA}_zxing_log.txt"
count_zbar=0
count_zxing=0

mkdir -p "$ANALYZE_ZBAR_DIR"
mkdir -p "$ANALYZE_ZXING_DIR"
mkdir -p "analyze/result_log/raw/zbar"
mkdir -p "analyze/result_log/raw/zxing"
mkdir -p "analyze/result_log/reg/zbar"
mkdir -p "analyze/result_log/reg/zxing"

for img in "$RENDER_DIR"/*.png; do
    filename=$(basename "$img" .png)
    result=$("$ZBARIMG_DIR" "$img" 2>&1)
    result_number=$(echo "$result" | grep -oP '(?<=:)\S+')
    echo "$result" > "$ANALYZE_ZBAR_DIR/${filename}.txt"
    echo "Processed: $filename -> $result_number" >> "$LOG_RAW_ZBAR"
    if [ "$result_number" = "$DATA" ]; then
        ((count_zbar++))
    fi
done

echo "$count_zbar" >> "$LOG_RAW_ZBAR"
./parse_log "$LOG_RAW_ZBAR" "$LOG_REG_ZBAR"
echo "Done zbar"

for img in "$RENDER_DIR"/*.png; do
    filename=$(basename "$img" .png)
    result=$("$ZXING_DIR" "$img" 2>&1)
    result_number=$(echo "$result" | grep -oP '(?<=Text:       )\S+' | tr -d '"')
    echo "$result" > "$ANALYZE_ZXING_DIR/${filename}.txt"
    echo "Processed: $filename -> $result_number" >> "$LOG_RAW_ZXING"
    if [ "$result_number" = "$DATA" ]; then
        ((count_zxing++))
    fi
done

echo "$count_zxing" >> "$LOG_RAW_ZXING"
./parse_log "$LOG_RAW_ZXING" "$LOG_REG_ZXING"
echo "Done zxing"

