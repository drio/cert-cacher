#!/bin/sh
#
set -eu

# run or print
rop() {
	local command=$*

	if [ $PRINT_ONLY -eq 1 ]; then
		echo "$command"
	else
		eval $command
	fi
}

tsCert() {
	rop "tailscale cert $DOMAIN"
}

saveCerts() {
	rop "curl --data-binary @./"$DOMAIN".cert $BASE_CATCHER_URL"
	rop "curl --data-binary @./"$DOMAIN".key $BASE_CATCHER_URL"
}

getCert() {
	rop "curl $BASE_CATCHER_URL/cert >$DOMAIN.cert"
	rop "curl $BASE_CATCHER_URL/key >$DOMAIN.key"
}

log() {
	local msg=$1
	printf "$(date)> $msg\n" >&2
}

main() {
	# How many days before the cert expires?
	response_body=$(mktemp)
	http_status=$(curl -s -o "$response_body" -w '%{http_code}' $BASE_CATCHER_URL/days)
	days_to_expire=$(cat $response_body)
	rm -f $response_body

	log "/days status=$http_status days_to_expire=$days_to_expire"

	if [ "$http_status" == "404" ]; then
		log "Cert not available in cacher. Requesting one and sending it to the cacher"
		tsCert
		saveCerts
	elif [ "$http_status" == "200" ]; then
		if [ $((days_to_expire > MIN_DAYS)) ]; then
			log "Cert cached and valid. Getting it from the catcher"
			getCert
		else
			log "Cert has expired in cacher. Requesting a new one and sending it to the cacher"
			tsCert
			saveCerts
		fi
	else
		echo "not expected http status: $http_status"
	fi
}

# If the cert has still MIN_DAYS before it expires we will use it.
# If not, we will issue a new cert via the tailscale client.
MIN_DAYS=30
BASE_CATCHER_URL="http://cert-cacher:9191"
DOMAIN=""
PRINT_ONLY=0

usage() {
	echo "Usage: $0 -d <domain> [-b <base_catcher_url> -m <min_days>]"
	echo "  -d <domain>          : full domain. Example: machine.tailnet.net"
	echo "  -b <base_catcher_url>: url to the catcher service ($BASE_CATCHER_URL)."
	echo "  -m <min_days>        : A cert needs min_days before it expires otherwise we will request a new one"
	echo "  -p                   : print cmds, do not execute them"
	echo ""
	echo "To download and execute: "
	echo "  curl -s http://cacher:9191/sh | sh -s -- -d foo.tailnet.net -b http://foo:1234 -m 50"
	exit 1
}

while getopts "b:d:m:p" opt; do
	case $opt in
	b) BASE_CATCHER_URL="$OPTARG" ;;
	d) DOMAIN="$OPTARG" ;;
	m) MIN_DAYS="$OPTARG" ;;
	p) PRINT_ONLY=1 ;;
	*) usage ;;
	esac
done

if [ ".$DOMAIN" == "." ]; then
	echo "Not domain provided. Bailing out."
	exit 1
fi

[ $PRINT_ONLY -eq 1 ] && log "-p enabled, printing cmds only "

main $BASE_CATCHER_URL $MIN_DAYS
