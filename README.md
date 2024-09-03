<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>


<!-- ABOUT THE PROJECT -->
## About The Project

This script does a very simple thing, it takes a path structure and converts it to products for a Woocommerce store.

It follows a very specific pattern and is not yet very flexible as it only supports one category, and a color property
derived from the name. 

Here's why:
* Nobody should be manually uploading and organizing pictures to upload, and then add description, SKU, variations etc...
* Photographers don't usually do any of the organization, however, they do have most of the information about the product, and variations available
* Using AI to create descriptions saves a ton of time :smile:


<!-- GETTING STARTED -->
## Getting Started

- build the executable with ```go build go build cmd/main/wcuploader.go```
- Set up the base directory to comply with the following structure
  - ```.base_dir/<category>/<Title> <ColorSKU>-<ColorName>-<Picture-perspective(front-Side..)>.jpg```
  - For example:
    - base_dir/kitchen/Can-Opener C1-PURPLE-FRONT.jpg
    - base_dir/tools/Drill Z2-Green-SideView.jpg
- call the executable created with the following flags:
```bash
./wcuploader -products "<base dir where directories with categories are"
-gptsecret "<openai-secret>" 
-shop-url "<wordpress base url>"
-woocommerce-key "<woocommerce-key>"
-woocommerce-secret "<woocommerce-secret>"
-wp-key "<wp app password>"
-user "<wp-user>"
```

### Prerequisites

Install Go: https://go.dev/learn/

### API Keys

#### How to obtain Woocommerce API Key

https://woocommerce.com/document/woocommerce-rest-api/

#### How to obtain Wordpress app password

https://wordpress.com/support/security/two-step-authentication/application-specific-passwords/

#### How to obtain OpenAI secret

https://platform.openai.com/api-keys