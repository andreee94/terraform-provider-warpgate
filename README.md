# terraform-provider-warpgate
An **unofficial** terraform provider for Warpgate 


## Notes

The client for the warpgate api is automatically generated with  
[oapi-codegen](github.com/deepmap/oapi-codegen/cmd/oapi-codegen) from the openapi specs at [https://raw.githubusercontent.com/warp-tech/warpgate/main/warpgate-web/src/admin/lib/openapi-schema.json](https://raw.githubusercontent.com/warp-tech/warpgate/main/warpgate-web/src/admin/lib/openapi-schema.json)

Unfortunately the specification uses `uint8` which are interpreted as base64 encoded string.

To solve this issue `uint8` are converted to `uint16` when generating the golang client.

To regenerate the client run the command:

```bash
make gen-warpgate
```


## Testing

<!-- To run acceptance testing is necessary to prepare the files for the warpgate docker container.

Unfortunately this step cannot be automated (at least I wasn't able), so a one time step is required:

```bash
sudo make gen-warpgate-setup
```

It will ask for `session recording` and a `password`.
Keep track of the password and update the `WARPGATE_PASSWORD` inside the `_scripts/testacc_setup.sh`, or just use the super bad `password` as password. 
This is just used for acceptance testing, and the container runs just for the time of the test. 

After the configuration files are generated, run the test with the command: -->

To perform the acceptance test run the command:

```bash
sudo make testacc
```

It will setup warpgate in unattended mode, run the container and perform the test against it.

> Docker is required
