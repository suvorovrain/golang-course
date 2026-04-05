package grpc

import(
	 "context"
    "fmt"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "github-info-system/api/generated"
)


type CollectorClient struct {
    client pb.CollectorServiceClient  
    conn   *grpc.ClientConn 
	
}



func NewCollectorClient(addr string) (*CollectorClient, error) { 
    conn, err := grpc.Dial(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithTimeout(5*time.Second),
    )
	
	
	
    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к Collector: %w", err)
    }

    return &CollectorClient{
        client: pb.NewCollectorServiceClient(conn),
		
		
		
        conn:   conn,
    }, nil
}



func (c *CollectorClient) GetRepoInfo(ctx context.Context, owner, repo string) (*pb.RepoResponse, error){
	req:= &pb.RepoRequest{
		Owner: owner,
		Repo: repo,
	}
	return c.client.GetRepoInfo(ctx, req) 
	
}

func (c *CollectorClient) Close() error{
	return c.conn.Close()
}


