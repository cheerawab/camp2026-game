#!/bin/sh
set -eu

url="${QUIZ_QUESTIONS_CSV_URL:-}"
if [ "$#" -ge 1 ] && [ -n "$1" ]; then
	url="$1"
fi

dest="${QUIZ_QUESTIONS_CSV_OUTPUT:-content/quiz_questions.csv}"
if [ "$#" -ge 2 ] && [ -n "$2" ]; then
	dest="$2"
fi

if [ -z "$url" ]; then
	echo "QUIZ_QUESTIONS_CSV_URL is required" >&2
	exit 1
fi

expected_header="id,prompt,choice_a,choice_b,choice_c,choice_d,correct_choice,explanation"

mkdir -p "$(dirname "$dest")"
tmp="${dest}.tmp"
rm -f "$tmp"
trap 'rm -f "$tmp"' EXIT HUP INT TERM

if command -v curl >/dev/null 2>&1; then
	curl -fsSL "$url" -o "$tmp"
elif command -v wget >/dev/null 2>&1; then
	wget -q -O "$tmp" "$url"
else
	echo "curl or wget is required to fetch quiz questions" >&2
	exit 1
fi

if [ ! -s "$tmp" ]; then
	echo "downloaded quiz questions CSV is empty" >&2
	exit 1
fi

header="$(sed -n '1{s/\r$//;p;}' "$tmp")"
if [ "$header" != "$expected_header" ]; then
	echo "unexpected quiz questions CSV header: $header" >&2
	echo "expected: $expected_header" >&2
	exit 1
fi

line_count="$(wc -l < "$tmp" | tr -d ' ')"
if [ "$line_count" -lt 11 ]; then
	echo "quiz questions CSV must contain at least 10 questions" >&2
	exit 1
fi

mv "$tmp" "$dest"
trap - EXIT HUP INT TERM
