# build
pushd server
make build_all

# upload
gsutil -m cp __dist/* gs://health-signal-deploy/server
popd

# download and restart
declare -a instances=(
  "--zone europe-west2-c healthsignal-1"
  )

for i in "${instances[@]}"
do
  echo "$i"
  gcloud compute ssh $i  -- bash -C refresh_server.sh
done
