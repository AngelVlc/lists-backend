# lists-backend

## To do

- Unify parse body methods and add body validations
- Expiration from ENV
- Add active field to user
- Refactor handler ServeHTTP

## Release image

```shell
docker build -t lists_release --target release .
```

```shell
docker run --rm -it -e PORT=5000 -e JWT_SECRET=the_jwt_secret -e MONGODB_URI=mongodb://host.docker.internal/listsDb -e MONGODB_URI_TEST=mongodb://host.docker.internal/listsTestDb lists_release
```

## Heroku

```shell
docker tag lists_release registry.heroku.com/app_name/worker
docker push registry.heroku.com/app_name/worker
heroku login
heroku container:login
heroku container:release --app app_name worker
```

## Terraform

- https://www.terraform.io/docs/providers/heroku/r/config.html
- https://devcenter.heroku.com/articles/using-terraform-with-heroku
