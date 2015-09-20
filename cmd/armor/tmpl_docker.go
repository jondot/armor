package main

const TEMPL_DOCKER = `
FROM tianon/true
COPY ./config /
COPY ./snapshot/linux_amd64/* /
EXPOSE 6060
CMD ["/{{.Product}}"]
`
