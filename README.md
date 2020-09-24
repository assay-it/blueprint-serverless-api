<p align="center">
  <img src="./blueprint.gif" height="240" />
  <h3 align="center">Blueprint: Serverless API</h3>
  <p align="center"><strong>How To Confirm Quality of Serverless API baked by DynamoDB, AWS Lambda and API Gateway</strong></p>

  <p align="center">
  </p>
</p>

--- 

Quality assurance of serverless applications is more complex than doing it for other runtimes. Engineering teams spend twice as much time maintaining testing environments and mocks of cloud dependencies instead of building a loyal relationship with their customers, [assay.it](https://assay.it) has you covered.

This example is inspired by the blog post [How To Confirm Quality of Serverless API baked by DynamoDB, AWS Lambda and API Gateway](https://assay.it/2020/09/24/confirm-quality-of-serverless-api/) and provides a reference implementation on classical blueprint of Serverless CRUD API.

The risk of occasional outage is always there due to regression in software quality despite the Serverless nature. Serverless only helps with operational and infrastructure management processes. However, engineers are still responsible for the quality of their products. Interfaces have to work as planned. It needs to serve api clients according to the value it promises and ensure great experience.

Safe testing of Serveless API requires automations, understanding of holistic approaches and proper isolation of test traffic from real one. A secret sauce of successful automation strategy is the code. Not only the Infrastructure as a Code but also expressing expected Behavior as a Code. Just write pure functional Golang instead of clicking through UI or maintaining endless XML, YAML or JSON documents.

## Getting Started

1. **Access to AWS Account** with permission to create/delete AWS resources is required to operate the blueprint

2. **Deploy** the blueprint. It built with declarative Infrastructure as a Code that describes the infrastructure using TypeScript and AWS CDK without explicit definition of commands to be performed.
```bash
npm -C cloud install
npm -C cloud run cdk -- deploy

...

 ✅  bookstore

Outputs:
bookstore.GatewayEndpoint4DF49EE0 = https://xxxxxxxxxx.execute-api.eu-west-1.amazonaws.com/api/

Stack ARN:
arn:aws:cloudformation:eu-west-1:000000000000:stack/bookstore/00000000-0000-0000-0000-000000000000
``` 

3. **Install** [assay command line](https://github.com/assay-it/assay)
```bash
go get github.com/assay-it/assay
```

4. **Sign up for [assay.it](https://assay.it)** with your GitHub developer account. Initially, the service requires only access to your public profile, public repositories and access to commit status.

5. **Fork [assay-it/blueprint-serverless-api](https://github.com/assay-it/blueprint-serverless-api)** to your own GitHub account and then add to the assay.it workspace. It implements a quality assessment suite for fictional bookstore Serverless API using [category pattern](https://assay.it/doc/core/category)
```go
func Create() assay.Arrow {
  book := Book{
    ID:    "book:hobbit",
    Title: "There and Back Again",
  }

  return http.Join(
    ø.POST("%s/books", sut),
    ø.ContentJSON(),
    ø.Send(&book),
    ƒ.Code(http.StatusCodeOK),
    ƒ.Recv(&book),
  ).Then(
    c.Value(&book.ID).String("book:hobbit"),
    c.Value(&book.Title).String("There and Back Again"),
  )
}
```

6. **Allow assay cli** to run quality assessment with [assay.it](https://assay.it) on your behalf, an access key is requires. Go to your profile settings at assay.it and generate a new personal access key.

7. **Run** quality assessment. The command uses a latest snapshot of Behavior as a Code suite from the repository to run the assessment against the given deployment of the service. 
```bash
assay run your-github/blueprint-serverless-api \
  --url https://xxxxxxxxxx.execute-api.eu-west-1.amazonaws.com/api \
  --key Z2l0aH...DBkdw
```

8. Results become visible at your [assay.it](https://assay.it) workspace.
![](https://assay.it/images/posts/2020-09-24-quality-assessment.png)


## Further Reading

Please continue to [the core](https://assay.it/doc/core) sections for details about Behavior as a Code development and see other [blueprints at GitHub](https://github.com/assay-it?q=blueprint).


## License

[![See LICENSE](https://img.shields.io/github/license/assay-it/blueprint-serverless-api.svg?style=for-the-badge)](LICENSE)
