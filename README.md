# Bitbucket Repository Cloner
It allows you to show a list of repositories or clone them to your local machine.

## Installation
To use this program, you need to have Go installed on your computer. If you don't have Go installed, you can download it from the official website.

## After installing Go, you can clone the program by running the following command:

```go
git clone git@github.com:indranandjha1993/bitBucketCloner.git
```

## Usage
To use the program, navigate to the directory where the program is installed and run the following command:

```go
go run main.go
```
The program will prompt you for your Bitbucket workspace name, username, and password. After entering this information, the program will ask you whether you want to show the list of repositories or clone them.

### Showing the list of repositories
To show the list of repositories, choose the s option. The program will prompt you for the number of repositories per page (default is 10) and the page number (default is 1). The program will then display the list of repositories.

### Cloning repositories
To clone repositories, choose the c option. The program will prompt you for the number of repositories per page (default is 10) and the page number (default is 1). The program will then clone all repositories to your local machine.

## License
This program is licensed under the MIT License. See the LICENSE file for details.