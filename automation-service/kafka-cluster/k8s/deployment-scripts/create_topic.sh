PARTICIPANT_ID=${PARTICIPANT_ID:-participant-id-test}

sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" ../create_topic_participant.template.yaml \
> create_topic_participant.$PARTICIPANT_ID.yaml

kubectl create -f ./create_topic_participant.$PARTICIPANT_ID.yaml
kubectl wait --timeout=60s --for=condition=complete job/$PARTICIPANT_ID-create-topic-participant
