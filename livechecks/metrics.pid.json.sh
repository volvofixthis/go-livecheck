timestamp=$(date +%s)
timestamp_new=$(($timestamp-5*60))
cat > ./output/metrics.2.json <<-EOF
{
    "gauge": {"client_connected": 1},
    "timer": {"last_ping": ${timestamp_new}}
}
EOF
cp ./output/metrics.2.json ./output/metrics.3.json
