#!/bin/bash

# Create root directory
mkdir -p cloudgallery

# Create all directories first
mkdir -p ./{cmd/{api,migrate,worker,tools},configs,internal/{shared/{events,valueobjects,errors,constants},domains/{auth/{domain,application,infrastructure/{repository,jwt,password}},gallery/{domain,application,infrastructure/{repository,query}},media/{domain,application,infrastructure/{processor,repository,cache}},sharing/{domain,application,infrastructure/{repository,token}},tenant/{domain,application,infrastructure/{repository,middleware}},subscription/{domain,application,infrastructure/{repository,payment_gateway,webhook}},payment/{domain,application,infrastructure/{repository,provider}},notification/{domain,application,infrastructure/{repository,provider,templates}},audit/{domain,application,infrastructure/{repository,exporter}}},infrastructure/{database/{postgres/{migrations,queries,sqlc_generated},redis},messaging/nats/dto,storage/{minio,s3},observability/{metrics,tracing,logging},cache,interfaces/{http/{handlers,middleware,dto},grpc/{proto,generated},cli},util},pkg/{errors,types,testing},api/{openapi,postman},scripts,deployments/{docker,kubernetes,docker-compose},docs/diagrams,test/{integration,contract,load},tools,web/{public,src}}

# Create all files (using touch to preserve existing content)
cd .

# cmd files
touch cmd/api/main.go
touch cmd/migrate/main.go
touch cmd/worker/{main.go,thumbnail_generator.go,email_sender.go,usage_calculator.go,cleanup_jobs.go}
touch cmd/tools/main.go

# configs files
touch configs/{config.local.yaml,config.staging.yaml,config.production.yaml,config.go}

# shared kernel files
touch internal/shared/events/{event_bus.go,nats_bus.go,domain_event.go}
touch internal/shared/valueobjects/{money.go,email.go,storage_quota.go}
touch internal/shared/errors/{domain_errors.go,application_errors.go}
touch internal/shared/constants/{plans.go,roles.go}

# auth domain
touch internal/domains/auth/domain/{user.go,session.go,domain_event.go,repository.go}
touch internal/domains/auth/application/{login.go,signup.go,logout.go,refresh_token.go,verify_email.go}
touch internal/domains/auth/infrastructure/repository/{postgres_user_repo.go,redis_session_repo.go}
touch internal/domains/auth/infrastructure/jwt/jwt_service.go
touch internal/domains/auth/infrastructure/password/bcrypt_hasher.go

# gallery domain
touch internal/domains/gallery/domain/{gallery.go,image.go,album.go,domain_event.go,repository.go}
touch internal/domains/gallery/application/{create_gallery.go,update_gallery.go,delete_gallery.go,upload_image.go,list_images.go,get_gallery_stats.go,organize_album.go}
touch internal/domains/gallery/infrastructure/repository/{postgres_gallery_repo.go,postgres_image_repo.go}
touch internal/domains/gallery/infrastructure/query/gallery_queries.sql

# media domain
touch internal/domains/media/domain/{image_variant.go,processing_job.go,metadata.go,domain_event.go,repository.go}
touch internal/domains/media/application/{process_image.go,generate_thumbnails.go,optimize_image.go,extract_metadata.go,convert_format.go}
touch internal/domains/media/infrastructure/processor/{libvips_processor.go,imagemagick_processor.go}
touch internal/domains/media/infrastructure/repository/postgres_processing_repo.go
touch internal/domains/media/infrastructure/cache/redis_metadata_cache.go

# sharing domain
touch internal/domains/sharing/domain/{share_link.go,share_permission.go,domain_event.go,repository.go}
touch internal/domains/sharing/application/{create_share_link.go,validate_share.go,revoke_share.go,get_shared_gallery.go,track_share_view.go}
touch internal/domains/sharing/infrastructure/repository/postgres_share_repo.go
touch internal/domains/sharing/infrastructure/token/secure_token_generator.go

# tenant domain
touch internal/domains/tenant/domain/{tenant.go,tenant_user.go,domain_event.go,repository.go}
touch internal/domains/tenant/application/{create_tenant.go,add_team_member.go,update_tenant_settings.go,get_tenant_usage.go,switch_tenant.go}
touch internal/domains/tenant/infrastructure/repository/postgres_tenant_repo.go
touch internal/domains/tenant/infrastructure/middleware/tenant_resolver.go

# subscription domain
touch internal/domains/subscription/domain/{subscription.go,plan.go,feature.go,invoice.go,domain_event.go,repository.go}
touch internal/domains/subscription/application/{subscribe_to_plan.go,cancel_subscription.go,change_plan.go,check_plan_limits.go,generate_invoice.go}
touch internal/domains/subscription/infrastructure/repository/{postgres_subscription_repo.go,postgres_plan_repo.go}
touch internal/domains/subscription/infrastructure/payment_gateway/{stripe_gateway.go,payment_interface.go}
touch internal/domains/subscription/infrastructure/webhook/stripe_webhook_handler.go

# payment domain
touch internal/domains/payment/domain/{payment.go,refund.go,payment_method.go,domain_event.go,repository.go}
touch internal/domains/payment/application/{process_payment.go,refund_payment.go,get_payment_history.go,handle_webhook.go}
touch internal/domains/payment/infrastructure/repository/postgres_payment_repo.go
touch internal/domains/payment/infrastructure/provider/stripe_client.go

# notification domain
touch internal/domains/notification/domain/{notification.go,email.go,template.go,domain_event.go,repository.go}
touch internal/domains/notification/application/{send_email.go,send_welcome_email.go,send_share_notification.go,send_invoice_email.go,queue_notification.go}
touch internal/domains/notification/infrastructure/repository/postgres_notification_repo.go
touch internal/domains/notification/infrastructure/provider/{sendgrid_provider.go,smtp_provider.go}
touch internal/domains/notification/infrastructure/templates/{welcome.html,share_invite.html,invoice.html}

# audit domain
touch internal/domains/audit/domain/{audit_log.go,usage_stat.go,domain_event.go,repository.go}
touch internal/domains/audit/application/{log_event.go,get_audit_trail.go,track_usage.go,generate_report.go}
touch internal/domains/audit/infrastructure/repository/{postgres_audit_repo.go,clickhouse_analytics_repo.go}
touch internal/domains/audit/infrastructure/exporter/prometheus_exporter.go

# infrastructure
touch internal/infrastructure/database/postgres/{connection.go,migrate.go}
touch internal/infrastructure/database/postgres/migrations/{000001_init_schema.up.sql,000001_init_schema.down.sql,000002_add_galleries.up.sql,000002_add_galleries.down.sql,000023_add_search_indexes.up.sql}
touch internal/infrastructure/database/postgres/queries/{auth.sql,gallery.sql,media.sql,sharing.sql,tenant.sql,subscription.sql,payment.sql,notification.sql,audit.sql}
touch internal/infrastructure/database/postgres/sqlc_generated/{db.go,models.go,queries.sql.go}
touch internal/infrastructure/database/redis/{connection.go,cache.go,session_store.go,rate_limiter.go}

# messaging
touch internal/infrastructure/messaging/nats/{connection.go,jetstream.go,publisher.go,subscriber.go,streams.go}
touch internal/infrastructure/messaging/dto/events.go

# storage
touch internal/infrastructure/storage/minio/{client.go,uploader.go,downloader.go,presigned_url.go}
touch internal/infrastructure/storage/s3/s3_storage.go
touch internal/infrastructure/storage/interface.go

# observability
touch internal/infrastructure/observability/metrics/{prometheus.go,counters.go,histograms.go}
touch internal/infrastructure/observability/tracing/{jaeger.go,otel.go}
touch internal/infrastructure/observability/logging/{logger.go,zap_logger.go}

# cache
touch internal/infrastructure/cache/{redis_cache.go,memory_cache.go}

# interfaces/http
touch internal/interfaces/http/handlers/{health_handler.go,auth_handler.go,gallery_handler.go,media_handler.go,sharing_handler.go,tenant_handler.go,subscription_handler.go,webhook_handler.go}
touch internal/interfaces/http/middleware/{auth.go,rate_limit.go,cors.go,logging.go,tenant.go,request_id.go,error_handler.go}
touch internal/interfaces/http/dto/{auth_requests.go,auth_responses.go,gallery_requests.go,gallery_responses.go,media_requests.go,media_responses.go,sharing_requests.go,subscription_requests.go,error_response.go}
touch internal/interfaces/http/{router.go,server.go}

# grpc
touch internal/interfaces/grpc/proto/{gallery.proto,media.proto,auth.proto}
touch internal/interfaces/grpc/server.go

# cli
touch internal/interfaces/cli/{commands.go,seed_data.go,cleanup.go}

# util
touch internal/util/{helpers.go,validator.go,crypto.go,pagination.go,date_utils.go}

# pkg
touch pkg/errors/errors.go
touch pkg/types/common.go
touch pkg/testing/{mocks.go,fixtures.go}

# api
touch api/openapi/{swagger.yaml,swagger.json}
touch api/postman/cloudgallery_collection.json

# scripts
touch scripts/{build.sh,deploy.sh,seed.sh,backup_db.sh,migrate.sh}
chmod +x scripts/*.sh

# deployments
touch deployments/docker/{Dockerfile,Dockerfile.worker,.dockerignore}
touch deployments/kubernetes/{deployment.yaml,service.yaml,ingress.yaml,configmap.yaml,secrets.yaml}
touch deployments/docker-compose/{docker-compose.yml,docker-compose.staging.yml,docker-compose.prod.yml}

# docs
touch docs/{architecture.md,api.md,deployment.md,development.md}
touch docs/diagrams/{architecture.png,event_flow.png}

# test
touch test/integration/{auth_test.go,gallery_test.go,api_test.go}
touch test/contract/stripe_contract_test.go
touch test/load/{upload_test.js,gallery_test.js}

# tools
touch tools/tools.go

# web
touch web/public/.gitkeep
touch web/src/.gitkeep

# root files
touch ../.env.example
touch ../.gitignore
touch ../.golangci.yml
touch ../Makefile
touch ../go.mod
touch ../go.sum
touch ../sqlc.yaml
touch ../.air.toml
touch ../README.md
touch ../LICENSE
touch ../CONTRIBUTING.md

echo "✅ Complete project structure created successfully!"
echo "📁 Location: $(pwd)/."
echo "📊 Total files created: $(find . -type f | wc -l)"
echo ""
echo "⚠️  IMPORTANT:"
echo "1. Your existing files have been preserved (no overwrites)"
echo "2. Empty placeholder files were created for the new structure"
echo "3. To migrate your existing code, run:"
echo "   find . -name '*.go' -not -path '././*' -type f"
echo ""
echo "🔄 Next steps:"
echo "1. Run 'go mod init' if needed in cloudgallery/"
echo "2. Copy your existing .go files to appropriate directories"
echo "3. Update import paths in your moved files"
echo "4. Run 'make build' to verify compilation"