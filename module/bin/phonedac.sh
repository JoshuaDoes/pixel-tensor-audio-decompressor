tinymix="$1"
mix () {
  $tinymix -D 0 $1 $2
  #if there is a right channel, set the offset
  if [ -n "$3" ]; then
    $tinymix -D 0 $(($1 + $3)) $2
  fi
}

pcm="${2:-817}" #817, max 913
amp="${3:-17}" #17, max 20
current="${4:-54}" #54, max 74

mix 5 $pcm 35
mix 6 $amp 35
mix 14 $current 35

# Clean up leftovers since the module now exits
rm -rf "/data/local/tmp/ptad"