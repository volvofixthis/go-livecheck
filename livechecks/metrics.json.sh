timestamp=$(date +%s)
timestamp=$(($timestamp-5*60))
cat <<-EOF
{
    "gauge": {"client_connected": 1},
    "timer": {"last_ping": ${timestamp}}
}
EOF

