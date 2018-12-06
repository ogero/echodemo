# ECHO Demo  

Basic demo of ECHO web project with some useful components plugged in for easy development and learning.    

## Some Features  

- High performance, minimalist Go web framework [ECHO](https://echo.labstack.com) 
- Users with roles and simple rbac  
- Basic settings  
- ORM (gorm)   
- Materializecss templates with snackbar/toasts  
- Encrypted (scrypt) cookie sessions  
- Tasks runner / scheduler  
- Locales enabled  
- Dependency Injection (goldi)  
- Mage for utils like binary build, binary embedding, assets packing, etc  
- Logrus logs  
- Automatic relaunch on sources changes  
- Single binary deploy with assets embedded  

## Requirements  
  
- GoLang 1.8.3+. GOPATH/bin folder must be in PATH env var.
- [dep](https://github.com/golang/dep/releases) as a dependency manager. Binary must reside in PATH env var.  
- [fileb0x](https://github.com/UnnoTed/fileb0x) as a virtual fs file embedding. Binary must reside in PATH env var.  
- [sassc binaries](https://github.com/eSlider/sassc-binaries/blob/develop/dist/sassc) as sassc preprocessor. Binary must reside in PATH env var.  
- [realize](https://gorealize.io/) as dev tooling. Binary must reside in PATH env var.  
- [mage](https://github.com/magefile/mage/releases) as build, dev and deploy tooling. Binary must reside in PATH env var.  

#### Usage  

1. Generate your `[echodemo.ini]`, config file options are available on `echodemo.ini.sample` file.    
2. Call the binary with a single optional parameter `[echodemo.ini]`.  

### Development  

1. Fetch dependencies: `dep ensure`  
2. On project root: `mage dev`  
3. Navigate [localhost:2332](http://localhost:2332), default user is `admin@localhost.com:1234`  
- Automatic restarts will be performed on source changes. See `.realize.yaml` config  

### Deployment  

To obtain a binary  
- On project root: `mage build`  
- Binaries will be placed in `dist` folder  
To deploy (task not runnable on windows)  
- Fetch dependencies (only when deploying for first time): `dep ensure`  
- On project root: `mage deploy`  
To make use of webhoooks (if configured on your version control platform and your server)  
- Make use of dist helper scripts (upstart, webhook config, sh script), note webhook config only whitelists bitbucket IPs  
- Do your dev stuff on master/other branch  
- Merge changes to production branch and push  
- Create a tag (following vN.N.N rule) on the merged commit of production branch and push (this will trigger a deploy)  