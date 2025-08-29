FROM public.ecr.aws/docker/library/golang AS build
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /binarytime-now ./cmd/lambda

FROM scratch
COPY --from=build /binarytime-now /binarytime-now
ENTRYPOINT ["/binarytime-now"]
