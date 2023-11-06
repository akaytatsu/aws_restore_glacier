this project is for restore a large quantity files in S3 from Glacier / Deep Archive / Glacier IR

# commands

## list all flles and write in file

```shell
    aws_restore_gaclier list_all --access_key AWS_ACCESS_KEY --secret_key AWS_SECRET_KEY --bucket BUCKET_NAME --region AWS_REGION_BUCKET --partial PREFIX_OPTIONAL
```

## list only files in glacier and write in file

```shell
    aws_restore_gaclier list --access_key AWS_ACCESS_KEY --secret_key AWS_SECRET_KEY --bucket BUCKET_NAME --region AWS_REGION_BUCKET --partial PREFIX_OPTIONAL
```

## restore files

```shell
    aws_restore_gaclier restore --access_key AWS_ACCESS_KEY --secret_key AWS_SECRET_KEY --bucket BUCKET_NAME --region AWS_REGION_BUCKET --partial PREFIX_OPTIONAL
```