fswatch -0 ./ | while read -d "" event; do
  echo "Changed: $event"
  go install ./...
done
