# stytch-go-magic-links

This is a lightweight [Stytch](https://stytch.com) + [Go](https://nextjs.org/) example app which demonstrates a quick Stytch implementation using our [Email Magic Links](https://stytch.com/docs/guides/magic-links/email-magic-links/api) product.

<p align="center"><img src="./public/example-app-image.png" alt="stytch" width="50%"/></p>

# Running locally

## Setting up Stytch
After signing up for Stytch, you'll need your Project's `project_id`, `secret`, and `public_token`. You can find these in the [API keys tab](https://stytch.com/dashboard/api-keys).

Once you've gathered these values, add them to a new .env.local file.
Example:

```bash
cp .env.template .env.local
# Replace your keys in new .env.local file
```

Next we'll configure the appropriate redirect URLs for your project, you'll set these magic link URLs for your project in the [Redirect URLs](https://stytch.com/dashboard/redirect-urls) section of your Dashboard. Add `http://localhost:3000/authenticate` as both a login and sign-up redirect URL. 

## Running the example app

Run `go run main.go`

Visit `http://localhost:3000` and login with your email.
Then check for the Stytch email and click the sign in button.
You should be signed in!
