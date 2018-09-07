# build
pushd client
make build_all

# upload
gsutil -m cp __dist/* gs://health-signal-deploy/client
popd

# download and restart
declare -a instances=(
  "--zone australia-southeast1-b healthsignal-sydney"
  "--zone europe-west2-c healthsignal-1"
  )

for i in "${instances[@]}"
do
  echo "$i"
  gcloud compute ssh $i  -- bash -C refresh_client.sh
done
