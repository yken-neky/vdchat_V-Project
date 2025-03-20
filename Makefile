# DEV

build-dev:
	docker build -t dev-video_chat_app -f containers/images/Dockerfile .

clean-dev:
	docker-compose -f containers/composes/dc.dev.yml down

run-dev:
	docker-compose -f containers/composes/dc.dev.yml up

logs-dev:
	docker-compose -f containers/composes/dc.dev.yml logs -f --tail 100

# PROD

stop-prod:
	docker-compose -f containers/composes/dc.prod.yml stop --timeout 120

clean-prod:
	docker-compose -f containers/composes/dc.prod.yml down --timeout 120

run-prod:
	docker-compose -f containers/composes/dc.prod.yml up -d --timeout 120

logs-prod:
	docker-compose -f containers/composes/dc.prod.yml logs -f --tail 100
