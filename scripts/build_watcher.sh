cmd="go install ./..."
$cmd
fswatch -0 ./ | while read -d "" event; do
  echo "Changed: $event"
  $cmd
done
