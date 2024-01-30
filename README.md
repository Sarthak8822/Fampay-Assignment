# FamPay Backend Assignment

## Project Goal

Create an API to fetch the latest YouTube videos for a given search query and store them in a database. Provide a paginated API to retrieve the stored videos in reverse chronological order.

## Basic Requirements

- Continuously fetch and store videos in the background with a 10-second interval.
- Implement a GET API to retrieve paginated video data sorted by publishing datetime.
- Ensure scalability and optimization.

## Bonus Points

- [x] Support multiple API keys for YouTube quota management.
- [ ] Optionally, create a dashboard to view stored videos with filters and sorting.

### Language & Framework

  `Golang`

### YouTube API and Golang Refernces

- YouTube Data v3 API: [YouTube Data v3 API](https://developers.google.com/youtube/v3/getting-started)
- Search API Reference: [Search API Reference](https://developers.google.com/youtube/v3/docs/search/list)
- Golang Youtube Package: [Discover Packages->google.golang.org/api->youtube->v3](https://pkg.go.dev/google.golang.org/api@v0.157.0/youtube/v3#Service.Search) 

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/Sarthak8822/Fampay-Assignment.git
   cd Fampay_Assignment
   
2. Create a .env file in the root of the project and add your configuration:

    ```env
    API_KEYS=your_api_key_1,your_api_key_2,your_api_key_3
    MONGO_URI=your_mongo_db_uri
    DATABASE_NAME=your_database_name
    COLLECTION_NAME=your_collection_name

  Replace your_api_key_1, your_mongo_db_uri, your_database_name, and your_collection_name with your actual API keys and MongoDB configuration.

3. Install dependencies:

    ```env
    go get -u ./...

4. Run the application::

    ```env
    go run main.go
    

## Endpoints

GET  

  `/api/paginated-videos`

  Retrieve paginated videos sorted in reverse chronological order of their publishing date-time.
  
  Query Parameters:
  
  - `page` `(optional, default: 1)`: `Page number`.
  - `pageSize` `(optional, default: 10)`: `Number of videos per page`.



## Background Fetching

  The project includes a continuous background fetching mechanism that fetches and stores YouTube videos with a specified query ("official" in this case). The background fetch is initiated in a goroutine when the server starts.

## My Resume
  [Sarthak_Modhe_Resume](https://drive.google.com/file/d/1y0GfuRRlVMOAiZ6gqAqj4-V4nWX0vY86/view?usp=drivesdk)

## Contribute
  If you'd like to contribute to this project, please follow the standard GitHub flow:

  1) Fork the repository.
  2) Create a new branch.
  3) Make your changes.
  4) Open a pull request.






