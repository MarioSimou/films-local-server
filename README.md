### Songs Local Server


### Pros and Cons

#### Disadvantages
- serverless-offline doesn't work correctly with binary data, which is not inline with API Gateway behavior. The difference between both environments make development and testing harder.

### Notes
- When an endpoint expects data in the form `multipart/form-data`, the binary file needs to be base64 encoded to prevent any potential encoding issues. This means, that the integration point (lambda function) has to decode the data, parse the binary file and upload to the bucket storage. Currently, `serverless-offline` doesn't support that behavior.
- If the endpoint needs to return binary data, the serverless template need to be updated and set the list of headers which should follow that behavior. Bear in mind, that this setup is applicable only on REST endpoints (v1). This is not available for HTTP API (v2) in AWS

```
provider:
    apiGateway:
        binaryMediaTypes:
            - 'multipart/form-data'
```

- If you get `illegal base64 data at input byte ...` while decoding a string, you most likely use the wrong decoder. [link](https://stackoverflow.com/questions/69753478/use-base64-stdencoding-or-base64-rawstdencoding-to-decode-base64-string-in-go). If you are wondering what is the difference between `StdEncoding` and `RawStdEncoding`, `RawStdEncoding` doesn't accept padding at the end of the string, while `StdEncoding` accepts.w
- API Gateway has a limit of 5MB of the payload, which prevent users to upload files of a greater size. If that's an issue for your implementation, have a look on pre-signed URLs.
- `serverless-apigw-binary` is intended for REST endpoint only.
- When an http event is configured with `cors: true`, it configures a preflight `OPTIONS` response for your endpoint, which is then responsible to return the necessary `Access-Control-Allow-*` headers.
- If you are working with binary data and the lambda function fails to extract the content-type header, please verify that is not a case-sensitiveness issue.  

### Resources
- [Working with binary media types for REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-payload-encodings.html)
- [API Gateway Conversions](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-payload-encodings-workflow.html)
- [Hit the 6MB Lambda payload limit?](https://theburningmonk.com/2020/04/hit-the-6mb-lambda-payload-limit-heres-what-you-can-do/)
- [What is the difference between StdEncoding and RawStdEncoding](https://pkg.go.dev/encoding/base64#:~:text=RawStdEncoding%20is%20the%20standard%20raw,StdEncoding%20but%20omits%20padding%20characters.&text=RawURLEncoding%20is%20the%20unpadded%20alternate%20base64%20encoding%20defined%20in%20RFC%204648.&text=StdEncoding%20is%20the%20standard%20base64%20encoding%2C%20as%20defined%20in%20RFC%204648.)