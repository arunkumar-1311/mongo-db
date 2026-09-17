# mongo-db

A Gin + MongoDB REST API for managing users, deployable to Kubernetes (minikube).

## Project Structure

- `main.go` — app entrypoint, connects to MongoDB and starts the router
- `dbconnection/` — MongoDB connection setup (`MONGO_URI`, `MONGO_DB` env vars)
- `models/` — data models (`User`, `Address`, `Department`, `Role`, `Order`, `LoginRequest`/`LoginResponse`)
- `auth/` — password hashing (bcrypt) and JWT generation/parsing (`JWT_SECRET` env var)
- `middleware/` — Gin middleware: `Authenticate` (JWT check), `RequireRole` (role-based authorization)
- `repository/` — MongoDB data access layer
- `service/` — business logic layer
- `handlers/` — HTTP handlers (Gin)
- `routers/` — route definitions, app listens on port `8000`
- `k8s/` — Kubernetes manifests
- `Dockerfile` — multi-stage build for the Go app

## Authentication & Authorization

- **Signup**: `POST /v1/api/user/create` — public. Creates a user with a bcrypt-hashed password (`password` in the request is never stored in plaintext).
- **Login**: `POST /v1/api/auth/login` — public. Verifies email/password and returns a signed JWT (24h expiry).
- **Protected routes**: e.g. `POST /v1/api/user/get/:id` requires `Authorization: Bearer <token>`. Enforced by `middleware.Authenticate()`.
- **Role-based authorization**: `middleware.RequireRole("admin", ...)` can be added to any route group to additionally restrict access by the caller's `role.level` (taken from the JWT claims). Apply it after `Authenticate()`.

Set `JWT_SECRET` to a long random value in any real environment — it defaults to a dev-only value if unset.

## API Endpoints

Base path: `/v1/api`

| Method | Path             | Auth required | Description        |
|--------|------------------|----------------|---------------------|
| POST   | `/user/create`   | No             | Create a user (signup) |
| POST   | `/auth/login`    | No             | Login, returns JWT  |
| POST   | `/user/get/:id`  | Yes (Bearer JWT) | Get a user by id  |

Example signup request body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SuperSecret123",
  "salary": 50000,
  "address": {
    "id": 1,
    "street": "MG Road",
    "city": "Bengaluru",
    "country": "India"
  },
  "department": {
    "name": "Engineering",
    "code": "ENG"
  },
  "role": {
    "title": "Developer",
    "level": "L2"
  }
}
```

Example login request body:
```json
{
  "email": "john@example.com",
  "password": "SuperSecret123"
}
```
Response:
```json
{ "token": "<jwt>" }
```

Call a protected route:
```bash
curl -X POST http://localhost:8000/v1/api/user/get/<id> \
  -H "Authorization: Bearer <jwt>"
```

## Local Development (without Kubernetes)

Requires a local MongoDB running on `localhost:27017` (or set env vars below).

```bash
# optional overrides, defaults shown
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB="local"
export JWT_SECRET="dev-secret-change-me"

go run main.go
```

App will be available at `http://localhost:8000`.

## Kubernetes Deployment (minikube)

### Prerequisites
- Docker
- minikube
- kubectl

### 1. Start minikube
```bash
minikube start
```

Check status any time:
```bash
minikube status
```

### 2. Build the Docker image

Build normally, then load it into minikube (most reliable method, works regardless of which Docker daemon built it):
```bash
docker build -t mongo-db-app:latest .
minikube image load mongo-db-app:latest
```

Verify the image is present inside minikube:
```bash
minikube image ls | grep mongo-db-app
```

### 3. Apply Kubernetes manifests

Apply in order:
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/mongo-secret.yaml
kubectl apply -f k8s/mongo-statefulset.yaml
kubectl apply -f k8s/app-configmap.yaml
kubectl apply -f k8s/app-secret.yaml
kubectl apply -f k8s/app-deployment.yaml
```

Or apply the whole folder at once (order not guaranteed, but fine after first-time setup):
```bash
kubectl apply -f k8s/
```

### 4. Verify the deployment

```bash
# all resources in the namespace
kubectl -n mongo-db-app get all

# pods only
kubectl -n mongo-db-app get pods

# watch pods come up live
kubectl -n mongo-db-app get pods -w

# deployment replica status (desired/ready/available)
kubectl -n mongo-db-app get deployment mongo-db-app

# count running app pods
kubectl -n mongo-db-app get pods -l app=mongo-db-app --no-headers | wc -l
```

Describe a pod to debug issues (e.g. `ImagePullBackOff`, `CrashLoopBackOff`):
```bash
kubectl -n mongo-db-app describe pod <pod-name>
```

View logs:
```bash
kubectl -n mongo-db-app logs <pod-name>
kubectl -n mongo-db-app logs -f <pod-name>          # follow/tail
kubectl -n mongo-db-app logs -l app=mongo-db-app     # all app pods
```

### 5. Access the API

Port-forward the service to your machine (keep this terminal open):
```bash
kubectl -n mongo-db-app port-forward svc/mongo-db-app 8000:80
```

Then call:
```
POST http://localhost:8000/v1/api/user/create      (signup, no auth)
POST http://localhost:8000/v1/api/auth/login        (login, no auth)
POST http://localhost:8000/v1/api/user/get/:id      (requires Authorization: Bearer <token>)
```

### 6. Restarting

Restart the app (e.g. after rebuilding the image or to pick up new pods):
```bash
kubectl -n mongo-db-app rollout restart deployment mongo-db-app
```

Restart MongoDB (data persists via PVC, safe to do):
```bash
kubectl -n mongo-db-app rollout restart statefulset mongo
```

Check rollout progress:
```bash
kubectl -n mongo-db-app rollout status deployment mongo-db-app
```

After a code change, rebuild and reload the image, then restart:
```bash
docker build -t mongo-db-app:latest .
minikube image load mongo-db-app:latest
kubectl -n mongo-db-app rollout restart deployment mongo-db-app
```

### 7. Scaling

```bash
kubectl -n mongo-db-app scale deployment mongo-db-app --replicas=3
```

### 8. Stopping / cleaning up

Delete just the app (keep MongoDB and its data):
```bash
kubectl delete -f k8s/app-deployment.yaml
```

Delete everything in the namespace:
```bash
kubectl delete namespace mongo-db-app
```

Stop minikube (keeps cluster state, faster to resume):
```bash
minikube stop
```

Delete minikube entirely (wipes all data including MongoDB PVC):
```bash
minikube delete
```

## MongoDB Data

Data is stored in a PersistentVolumeClaim (`mongo-data`) mounted at `/data/db` in the `mongo-0` pod, defined in `k8s/mongo-statefulset.yaml`. It persists across pod/deployment restarts, but is wiped if you run `minikube delete`.

Check the PVC:
```bash
kubectl -n mongo-db-app get pvc
kubectl -n mongo-db-app get pv
```

Connect to MongoDB directly to inspect data:
```bash
kubectl -n mongo-db-app exec -it mongo-0 -- mongosh -u root -p changeme --authenticationDatabase admin
```

Inside the `mongosh` shell:
```js
use appdb
db.users.find()
```

(Default credentials are defined in `k8s/mongo-secret.yaml` — change `MONGO_INITDB_ROOT_PASSWORD` there before any real/shared use. Likewise change `JWT_SECRET` in `k8s/app-secret.yaml`.)

## Kubernetes Dashboard (UI)

minikube ships a web UI for viewing pods, deployments, logs, etc.:
```bash
minikube dashboard
```
This opens a browser and blocks the terminal until `Ctrl+C`. Use the namespace dropdown to switch to `mongo-db-app`.

To just print the URL instead of opening a browser:
```bash
minikube dashboard --url
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| `ImagePullBackOff` | Image wasn't built into minikube's environment | `minikube image load mongo-db-app:latest`, then `kubectl rollout restart deployment mongo-db-app` |
| `ECONNREFUSED` calling `localhost:8000` | `port-forward` not running | Run the port-forward command and keep it open |
| 404 on API call | Missing `/v1/api` prefix | Use `http://localhost:8000/v1/api/user/create` |
| Pod stuck `Pending` | Not enough resources / PVC not bound | `kubectl -n mongo-db-app describe pod <name>` to see the reason |
