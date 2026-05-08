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

mkdir -p "$ANALYZE_ZBAR_DIR"
mkdir -p "$ANALYZE_ZXING_DIR"
mkdir -p "analyze/result_log/raw/zbar"
mkdir -p "analyze/result_log/raw/zxing"
mkdir -p "analyze/result_log/reg/zbar"
mkdir -p "analyze/result_log/reg/zxing"

process_zbar() {
    local img="$1"
    local filename=$(basename "$img" .png)
    local result=$("$ZBARIMG_DIR" "$img" 2>&1)
    local result_number=$(echo "$result" | grep -oP '(?<=:)\S+')
    echo "$result" > "$ANALYZE_ZBAR_DIR/${filename}.txt"
    echo "Processed: $filename -> $result_number" >> "$LOG_RAW_ZBAR"
    echo "$ANALYZE_ZBAR_DIR/${filename}.txt"
}

export -f process_zbar
export ZBARIMG_DIR ANALYZE_ZBAR_DIR LOG_RAW_ZBAR

ls "$RENDER_DIR"/*.png | xargs -P 4 -I {} bash -c 'process_zbar "$@"' _ {}

if [ ! -f "./parse_log" ]; then
    go build -o parse_log parse_log.go 
fi

echo "analyze zbar log"
echo "" >> "$LOG_RAW_ZBAR"
./parse_log "$LOG_RAW_ZBAR" "$LOG_REG_ZBAR"
echo "Done zbar"

process_zxing() {
    local img="$1"
    local filename=$(basename "$img" .png)
    local result=$("$ZXING_DIR" "$img" 2>&1)
    local result_number=$(echo "$result" | grep -oP '(?<=Text:       )\S+' | tr -d '"')
    echo "$result" > "$ANALYZE_ZXING_DIR/${filename}.txt"
    echo "Processed: $filename -> $result_number" >> "$LOG_RAW_ZXING"
    echo "$ANALYZE_ZXING_DIR/${filename}.txt"
}

export -f process_zxing
export ZXING_DIR ANALYZE_ZXING_DIR LOG_RAW_ZXING

ls "$RENDER_DIR"/*.png | xargs -P 4 -I {} bash -c 'process_zxing "$@"' _ {}

echo "analyze zxing log"
echo "" >> "$LOG_RAW_ZXING"
./parse_log "$LOG_RAW_ZXING" "$LOG_REG_ZXING"
echo "Done zxing"