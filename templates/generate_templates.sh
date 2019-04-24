#!/bin/bash

# overwrite settings from default file
[ -e /etc/sysconfig/resistor ] && . /etc/sysconfig/resistor


RESISTOR_SERVER=${RESISTOR_URL:-"http://localhost:6090"}
RESISTOR_USERNAME=${RESISTOR_USER:-"adm1"} 
RESISTOR_PASSWORD=${RESISTOR_PASS:-"adm1pass"}


WORK_DIR=${RESIST_HOME:-$PWD}



OUTPUT_DIR="${WORK_DIR}/generated_tpls"
OUTPUT_DIR_TH="${WORK_DIR}/generated_tpls/th"
OUTPUT_DIR_AT="${WORK_DIR}/generated_tpls/at"
OUTPUT_DIR_RT="${WORK_DIR}/generated_tpls/rt"
OUTPUT_DIR_DEADMAN="${WORK_DIR}/generated_tpls/deadman"

. ${WORK_DIR}/available_vars.tpl

COOKFILE=$(mktemp)
TPL_JSON_FILE=$(mktemp)

function json_escape() {
  cat $1 | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
}

function log_to_resistor() {
  curl -s -X POST -F "username=$1" -F "password=$2" -c  $COOKFILE ${RESISTOR_SERVER}/login -o /dev/null
}

# put_to_resistor ID TRIGERTYPE CRITDIRECTION TTYPE STATFUNC THRESHOLDTYPE TICK_FILE DESCRIPTION
function put_to_resistor() {
# check if exist
HTTP_OUT=`curl  -sL -w "%{http_code}\\n" -b $COOKFILE  ${RESISTOR_SERVER}/api/cfg/template/${1} -o /dev/null`

if [ "$HTTP_OUT" == "200" ];
then
# ID exist , we should update
  METHOD="PUT"
  ID="$1"
else
  METHOD="POST"
  ID=""
fi




case $2 in
  "DEADMAN")
cat << EOF > $TPL_JSON_FILE
{
  "ID": "$1",
  "TriggerType": "$2",
  "TplData": $(json_escape $3),
  "Description": "$4"
}
EOF
  ;;
  "TREND")

  case $6 in
    "AT")
    TRENDTYPE="absolute"
    ;;
    "RT")
    TRENDTYPE="relative"
    ;;
  esac 
  case $8 in
    "P")
    TREND_SIGN="positive"
    ;;
    "N")
    TREND_SIGN="negative"
    ;;
  esac 

cat << EOF > $TPL_JSON_FILE
{
  "ID": "$1",
  "TriggerType": "$2",
  "CritDirection": "$3",
	"FieldType" : "$4",
  "StatFunc": "$5",
  "TrendType": "$TRENDTYPE",
  "TrendSign" : "$TREND_SIGN",
  "TplData": $(json_escape $7),
  "Description": "$9"
}
EOF
  ;;
  "THRESHOLD")

cat << EOF > $TPL_JSON_FILE
{
  "ID": "$1",
  "TriggerType": "$2",
  "CritDirection": "$3",
	"FieldType" : "$4",
  "StatFunc": "$5",
  "TplData": $(json_escape $7),
  "Description": "$8"
}
EOF
  ;;
esac

curl -s -X $METHOD -b $COOKFILE -d @$TPL_JSON_FILE ${RESISTOR_SERVER}/api/cfg/template/${ID} --header "Content-Type: application/json" -o /dev/null
OUT=$?
if [ "$OUT" -eq 0 ] 
then
  echo "Loaded OK"
else
  echo "ERROR CODE ${OUT}"
fi 
}



# Init Login

log_to_resistor  ${RESISTOR_USERNAME} ${RESISTOR_PASSWORD}

#LOAD VARS



#Create TPL structure

if [ -d "$OUTPUT_DIR" ]
then
	rm -rf $OUTPUT_DIR
fi

mkdir -p $OUTPUT_DIR_TH
mkdir -p $OUTPUT_DIR_AT
mkdir -p $OUTPUT_DIR_RT
mkdir -p $OUTPUT_DIR_DEADMAN

## Iterate over all functions:

for ttype in "${TEMPLATE_TYPE[@]}"
do
	if [ "$ttype" == "COUNTER" ]
	then
		derivnode=${DERIV_NODE}
		unitdata="\\n//Unit var to derivative node\\nvar UNIT_DERIV = 1m\\n"
	else
		derivnode=''
		unitdata=''
	fi
	printf "[%s]\n" "$ttype"
	for function in "${AVAILABLE_FUNCTIONS[@]}"
	do
		if [ "$function" == "movingAverage" ]
		then
			window=''
			extravar="\\n//ExtraData to configure the function\\nvar EXTRA_DATA = 5\\n"
			extranode=', EXTRA_DATA'
		elif [ "$function" ==  "percentile" ]
		then
			window=${WINDOW_NODE}
			extravar="\\n//ExtraData to configure the function\\nvar EXTRA_DATA = 95.0\\n"
			extranode=', EXTRA_DATA'
		else
			window=${WINDOW_NODE}
			extravar=''
			extranode=''
		fi
		
		printf "\t[%s]\n" "$function"

		## Generate UA templates
		TYPE="THRESHOLD"
		SUBTYPE="TH"
		OUTPUT_TPL_DIR=$OUTPUT_DIR_TH

		for direction in "${AVAILABLE_TH_DIRECTIONS[@]}"
		do
			if [ "$direction" == ">" ] 
			then
				DIRECTION="AC"
			else
				DIRECTION="DC"
			fi

			FUNCTION=`echo $function |awk '{print toupper($0)}'`
			TPL_ID=${TYPE}_2EX_${DIRECTION}_${SUBTYPE}_${ttype}_F${FUNCTION}
			TPL_OUTPUT=${OUTPUT_TPL_DIR}/${TPL_ID}.tpl

			sed	-e "s|@EXTRADATA@|$extravar|g" \
					-e "s|@EXTRANODE@|$extranode|g" \
					-e "s|@UNITDATA@|$unitdata|g" \
					-e "s|@DERIVNODE@|$derivnode|g" \
					-e "s|@WINDOW@|$window|g" \
					-e "s|@FUNCTION@|$function|g" \
					-e "s|@DIRECTION@|$direction|g" ${WORK_DIR}/tpl/kapacitor_TH_template.tpl > $TPL_OUTPUT
			http_out=`put_to_resistor ${TPL_ID} ${TYPE} ${DIRECTION} ${ttype} ${FUNCTION} ${SUBTYPE} ${TPL_OUTPUT} ""`
			printf "\t\t|-[%s][%s][%s][%s]\tGenerated %s [HTTP %s]\n" "$ttype" "$FUNCTION" "$TYPE" "$SUBTYPE" "$TPL_OUTPUT" "$http_out"
		done


		for trend_type in "${TREND_TYPE[@]}"
		do
			if [ "$trend_type" == "" ]
			then
				TYPE="TREND"
				SUBTYPE="AT"
				OUTPUT_TPL_DIR=$OUTPUT_DIR_AT
			else
				TYPE="TREND"
				SUBTYPE="RT"
				OUTPUT_TPL_DIR=$OUTPUT_DIR_RT
			fi

			## Generate TXX_CC template
			for direction in "${AVAILABLE_TXX_AC_DIRECTIONS[@]}"
			do
				if [ "$direction" == ">" ] 
				then
					VALUE_TYPE="P"
				else
					VALUE_TYPE="N"
				fi
				DIRECTION="AC"
				FUNCTION=`echo $function |awk '{print toupper($0)}'`
				TPL_ID=${TYPE}_2EX_${DIRECTION}_${SUBTYPE}${VALUE_TYPE}_${ttype}_F${FUNCTION}
				TPL_OUTPUT=${OUTPUT_TPL_DIR}/${TPL_ID}.tpl

				sed	-e "s|@EXTRADATA@|$extravar|g" \
						-e "s|@EXTRANODE@|$extranode|g" \
						-e "s|@WINDOW@|$window|g" \
						-e "s|@UNITDATA@|$unitdata|g" \
						-e "s|@DERIVNODE@|$derivnode|g" \
						-e "s|@FUNCTION@|$function|g" \
						-e "s|@TREND_TYPE@|$trend_type|g" \
						-e "s|@DIRECTION@|$direction|g" ${WORK_DIR}/tpl/kapacitor_TXX_AC_template.tpl > ${TPL_OUTPUT}
				http_out=`put_to_resistor ${TPL_ID} ${TYPE} ${DIRECTION} ${ttype} ${FUNCTION} ${SUBTYPE} ${TPL_OUTPUT}  ${VALUE_TYPE} ""`
				printf "\t\t|-[%s][%s][%s][%s]\t\tGenerated %s [HTTP %s]\n" "$ttype" "$FUNCTION" "$TYPE" "$SUBTYPE" "$TPL_OUTPUT" "$http_out"
			done

			## Generate TXX_CD template
			for direction in "${AVAILABLE_TXX_DC_DIRECTIONS[@]}"
			do
				if [ "$direction" == "<" ] 
				then
					VALUE_TYPE="P"
					direction0=${AVAILABLE_TXP_DC_DIRECTIONS0}
				else
					VALUE_TYPE="N"
					direction0=${AVAILABLE_TXN_DC_DIRECTIONS0}
				fi
				DIRECTION="DC"
				FUNCTION=`echo $function |awk '{print toupper($0)}'`
				TPL_ID=${TYPE}_2EX_${DIRECTION}_${SUBTYPE}${VALUE_TYPE}_${ttype}_F${FUNCTION}
				TPL_OUTPUT=${OUTPUT_TPL_DIR}/${TPL_ID}.tpl

				sed	-e "s|@EXTRADATA@|$extravar|g" \
						-e "s|@EXTRANODE@|$extranode|g" \
						-e "s|@WINDOW@|$window|g" \
						-e "s|@UNITDATA@|$unitdata|g" \
						-e "s|@DERIVNODE@|$derivnode|g" \
						-e "s|@FUNCTION@|$function|g" \
						-e "s|@TREND_TYPE@|$trend_type|g" \
						-e "s|@DIRECTION0@|$direction0|g" \
						-e "s|@DIRECTION@|$direction|g" ${WORK_DIR}/tpl/kapacitor_TXX_DC_template.tpl > ${TPL_OUTPUT}
				http_out=`put_to_resistor ${TPL_ID} ${TYPE} ${DIRECTION} ${ttype} ${FUNCTION} ${SUBTYPE} ${TPL_OUTPUT} ""`
				printf "\t\t|-[%s][%s][%s][%s]\t\tGenerated %s [HTTP %s]\n" "$ttype" "$FUNCTION" "$TYPE" "$SUBTYPE" "$TPL_OUTPUT" "$http_out"
			done
		done
	done
done

## Generates the DEADMAN alert

OUTPUT_TPL_DIR=$OUTPUT_DIR_DEADMAN

TPL_OUTPUT=${OUTPUT_TPL_DIR}/DEADMAN.tpl

cat ${WORK_DIR}/tpl/kapacitor_DEADMAN_template.tpl > ${TPL_OUTPUT}
http_out=`put_to_resistor "DEADMAN" "DEADMAN" ${TPL_OUTPUT} ""`


printf "[DEADMAN]\n"

printf "\t\t|-[DEADMAN][--]\t\tGenerated %s [HTTP %s]\n" "$TPL_OUTPUT" "$http_out"

rm -f $COOKFILE
rm -f $TPL_JSON_FILE
