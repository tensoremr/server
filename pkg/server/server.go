/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "net/http/pprof"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/robfig/cron/v3"
	"github.com/tensoremr/server/pkg/auth"
	"github.com/tensoremr/server/pkg/conf"
	"github.com/tensoremr/server/pkg/controller"
	"github.com/tensoremr/server/pkg/graphql/graph"
	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	"gorm.io/gorm"
)

// Server ...
type Server struct {
	Gin           *gin.Engine
	Config        *conf.Configuration
	DB            *gorm.DB
	ACLEnforcer   *casbin.Enforcer
	TestDB        *gorm.DB      // Database connection
	ModelRegistry *models.Model // Model registry for migration
}

// NewServer will create a new instance of the application
func NewServer() *Server {
	server := &Server{}

	server.ModelRegistry = models.NewModel()
	server.NewEnforcer()

	if err := server.ModelRegistry.OpenPostgres(); err != nil {
		log.Fatalf("gorm: could not connect to db %q", err)
	}

	//injectDb(server.ModelRegistry.DB)

	server.DB = server.ModelRegistry.DB

	server.ModelRegistry.RegisterAllModels()
	server.ModelRegistry.AutoMigrateAll()

	server.AddSearchIndex()
	server.SeedData()
	server.RegisterJobs()

	server.Gin = server.NewRouter()

	return server
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RegisterJobs ...
func (s *Server) RegisterJobs() {
	patientQueueRepository := repository.ProvidePatientQueueRepository(s.DB)

	c := cron.New()
	c.AddFunc("@hourly", func() {
		if err := patientQueueRepository.ClearExpired(); err != nil {
			fmt.Println(err)
		}
	})
	c.Start()
}

// Defining the Graphql handler
func graphqlHandler(server *Server, h *handler.Server) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// NewRouter ...
func (s *Server) NewRouter() *gin.Engine {
	AllergyRepository := repository.ProvideAllergyRepository(s.DB)
	AmendmentRepository := repository.ProvideAmendmentRepository(s.DB)
	AppointmentQueueRepository := repository.ProvideAppointmentQueueRepository(s.DB)
	AppointmentStatusRepository := repository.ProvideAppointmentStatusRepository(s.DB)
	AppointmentRepository := repository.ProvideAppointmentRepository(s.DB, AppointmentStatusRepository)
	AutoRefractionRepository := repository.ProvideAutoRefractionRepository(s.DB)
	BillingRepository := repository.ProvideBillingRepository(s.DB)
	ChatDeleteRepository := repository.ProvideChatDeleteRepository(s.DB)
	ChatMemberRepository := repository.ProvideChatMemberRepository(s.DB)
	ChatMessageRepository := repository.ProvideChatMessageRepository(s.DB)
	ChatMuteRepository := repository.ProvideChatMuteRepository(s.DB)
	ChatUnreadRepository := repository.ProvideChatUnreadRepository(s.DB)
	ChatRepository := repository.ProvideChatRepository(s.DB)
	ChiefComplaintRepository := repository.ProvideChiefComplaintRepository(s.DB)
	CoverTestRepository := repository.ProvideCoverTestRepository(s.DB)
	DiagnosticProcedureRepository := repository.ProvideDiagnosticProcedureRepository(s.DB)
	DiagnosticProcedureOrderRepository := repository.ProvideDiagnosticProcedureOrderRepository(s.DB)
	DiagnosticProcedureTypeRepository := repository.ProvideDiagnosticProcedureTypeRepository(s.DB)
	ExamCategoryRepository := repository.ProvideExamCategoryRepository(s.DB)
	ExternalExamRepository := repository.ProvideExternalExamRepository(s.DB)
	ExamFindingRepository := repository.ProvideExamFindingRepository(s.DB)
	EyewearPrescriptionOrderRepository := repository.ProvideEyewearPrescriptionOrderRepository(s.DB)
	EyewearPrescriptionRepository := repository.ProvideEyewearPrescriptionRepository(s.DB)
	EyewearShopRepository := repository.ProvideEyewearShopRepository(s.DB)
	FamilyIllnessRepository := repository.ProvideFamilyIllnessRepository(s.DB)
	FavoriteChiefComplaintRepository := repository.ProvideFavoriteChiefComplaintRepository(s.DB)
	ChiefComplaintTypeRepository := repository.ProvideChiefComplaintTypeRepository(s.DB, FavoriteChiefComplaintRepository)
	FavoriteDiagnosisRepository := repository.ProvideFavoriteDiagnosisRepository(s.DB)
	DiagnosisRepository := repository.ProvideDiagnosisRepository(s.DB, FavoriteDiagnosisRepository)
	FavoriteMedicationRepository := repository.ProvideFavoriteMedicationRepository(s.DB)
	FileRepository := repository.ProvideFileRepository(s.DB)
	FollowUpOrderRepository := repository.ProvideFollowUpOrderRepository(s.DB)
	FollowUpRepository := repository.ProvideFollowUpRepository(s.DB)
	FunduscopyRepository := repository.ProvideFunduscopyRepository(s.DB)
	HpiComponentTypeRepository := repository.ProvideHpiComponentTypeRepository(s.DB)
	HpiComponentRepository := repository.ProvideHpiComponentRepository(s.DB)
	IopRepository := repository.ProvideIopRepository(s.DB)
	LabOrderRepository := repository.ProvideLabOrderRepository(s.DB)
	LabTypeRepository := repository.ProvideLabTypeRepository(s.DB)
	LabRepository := repository.ProvideLabRepository(s.DB)
	LifestyleTypeRepository := repository.ProvideLifestyleTypeRepository(s.DB)
	LifestyleRepository := repository.ProvideLifestyleRepository(s.DB)
	MedicalPrescriptionOrderRepository := repository.ProvideMedicalPrescriptionOrderRepository(s.DB)
	MedicalPrescriptionRepository := repository.ProvideMedicalPrescriptionRepository(s.DB)
	OcularMotilityRepository := repository.ProvideOcularMotilityRepository(s.DB)
	OpthalmologyExamRepository := repository.ProvideOpthalmologyExamRepository(s.DB)
	OpticDiscRepository := repository.ProvideOpticDiscRepository(s.DB)
	OrganizationDetailsRepository := repository.ProvideOrganizationDetailsRepository(s.DB)
	PastHospitalizationRepository := repository.ProvidePastHospitalizationRepository(s.DB)
	PastIllnessTypeRepository := repository.ProvidePastIllnessTypeRepository(s.DB)
	PastIllnessRepository := repository.ProvidePastIllnessRepository(s.DB)
	PastInjuryRepository := repository.ProvidePastInjuryRepository(s.DB)
	PastOptSurgeryRepository := repository.ProvidePastOptSurgeryRepository(s.DB)
	PastSurgeryRepository := repository.ProvidePastSurgeryRepository(s.DB)
	PatientChartRepository := repository.ProvidePatientChartRepository(s.DB)
	PatientDiagnosisRepository := repository.ProvidePatientDiagnosisRepository(s.DB)
	PatientEncounterLimitRepository := repository.ProvidePatientEncounterLimitRepository(s.DB)
	PatientHistoryRepository := repository.ProvidePatientHistoryRepository(s.DB)
	PatientQueueRepository := repository.ProvidePatientQueueRepository(s.DB)
	PatientRepository := repository.ProvidePatientRepository(s.DB)
	PaymentWaiverRepository := repository.ProvidePaymentWaiverRepository(s.DB)
	PaymentRepository := repository.ProvidePaymentRepository(s.DB)
	PharmacyRepository := repository.ProvidePharmacyRepository(s.DB)
	PhysicalExamFindingRepository := repository.ProvidePhysicalExamFindingRepository(s.DB)
	PupilsRepository := repository.ProvidePupilsRepository(s.DB)
	QueueDestinationRepository := repository.ProvideQueueDestinationRepository(s.DB)
	QueueSubscriptionRepository := repository.ProvideQueueSubscriptionRepository(s.DB)
	ReferralOrderRepository := repository.ProvideReferralOrderRepository(s.DB)
	ReferralRepository := repository.ProvideReferralRepository(s.DB)
	ReviewOfSystemRepository := repository.ProvideReviewOfSystemRepository(s.DB)
	RoomRepository := repository.ProvideRoomRepository(s.DB)
	SlitLampExamRepository := repository.ProvideSlitLampExamRepository(s.DB)
	SupplyRepository := repository.ProvideSupplyRepository(s.DB)
	SurgicalOrderRepository := repository.ProvideSurgicalOrderRepository(s.DB)
	SurgicalProcedureTypeRepository := repository.ProvideSurgicalProcedureTypeRepository(s.DB)
	SurgicalProcedureRepository := repository.ProvideSurgicalProcedureRepository(s.DB)
	SystemSymptomRepository := repository.ProvideSystemSymptomRepository(s.DB)
	SystemRepository := repository.ProvideSystemRepository(s.DB)
	TreatmentOrderRepository := repository.ProvideTreatmentOrderRepository(s.DB)
	TreatmentTypeRepository := repository.ProvideTreatmentTypeRepository(s.DB)
	TreatmentRepository := repository.ProvideTreatmentRepository(s.DB)
	UserTypeRepository := repository.ProvideUserTypeRepository(s.DB)
	UserRepository := repository.ProvideUserRepository(s.DB, UserTypeRepository)
	VisitTypeRepository := repository.ProvideVisitTypeRepository(s.DB)
	VisualAcuityRepository := repository.ProvideVisualAcuityRepository(s.DB)
	VitalSignsRepository := repository.ProvideVitalSignsRepository(s.DB)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Config:                             s.Config,
		AccessControl:                      s.ACLEnforcer,
		AllergyRepository:                  AllergyRepository,
		AmendmentRepository:                AmendmentRepository,
		AppointmentQueueRepository:         AppointmentQueueRepository,
		AppointmentStatusRepository:        AppointmentStatusRepository,
		AppointmentRepository:              AppointmentRepository,
		AutoRefractionRepository:           AutoRefractionRepository,
		BillingRepository:                  BillingRepository,
		ChatDeleteRepository:               ChatDeleteRepository,
		ChatMemberRepository:               ChatMemberRepository,
		ChatMessageRepository:              ChatMessageRepository,
		ChatMuteRepository:                 ChatMuteRepository,
		ChatUnreadRepository:               ChatUnreadRepository,
		ChatRepository:                     ChatRepository,
		ChiefComplaintTypeRepository:       ChiefComplaintTypeRepository,
		ChiefComplaintRepository:           ChiefComplaintRepository,
		CoverTestRepository:                CoverTestRepository,
		DiagnosisRepository:                DiagnosisRepository,
		DiagnosticProcedureRepository:      DiagnosticProcedureRepository,
		DiagnosticProcedureOrderRepository: DiagnosticProcedureOrderRepository,
		DiagnosticProcedureTypeRepository:  DiagnosticProcedureTypeRepository,
		ExamCategoryRepository:             ExamCategoryRepository,
		ExternalExamRepository:             ExternalExamRepository,
		ExamFindingRepository:              ExamFindingRepository,
		EyewearPrescriptionOrderRepository: EyewearPrescriptionOrderRepository,
		EyewearPrescriptionRepository:      EyewearPrescriptionRepository,
		EyewearShopRepository:              EyewearShopRepository,
		FamilyIllnessRepository:            FamilyIllnessRepository,
		FavoriteChiefComplaintRepository:   FavoriteChiefComplaintRepository,
		FavoriteDiagnosisRepository:        FavoriteDiagnosisRepository,
		FavoriteMedicationRepository:       FavoriteMedicationRepository,
		FileRepository:                     FileRepository,
		FollowUpOrderRepository:            FollowUpOrderRepository,
		FollowUpRepository:                 FollowUpRepository,
		FunduscopyRepository:               FunduscopyRepository,
		HpiComponentTypeRepository:         HpiComponentTypeRepository,
		HpiComponentRepository:             HpiComponentRepository,
		IopRepository:                      IopRepository,
		LabOrderRepository:                 LabOrderRepository,
		LabTypeRepository:                  LabTypeRepository,
		LabRepository:                      LabRepository,
		LifestyleTypeRepository:            LifestyleTypeRepository,
		LifestyleRepository:                LifestyleRepository,
		MedicalPrescriptionOrderRepository: MedicalPrescriptionOrderRepository,
		MedicalPrescriptionRepository:      MedicalPrescriptionRepository,
		OcularMotilityRepository:           OcularMotilityRepository,
		OpthalmologyExamRepository:         OpthalmologyExamRepository,
		OpticDiscRepository:                OpticDiscRepository,
		OrganizationDetailsRepository:      OrganizationDetailsRepository,
		PastHospitalizationRepository:      PastHospitalizationRepository,
		PastIllnessTypeRepository:          PastIllnessTypeRepository,
		PastIllnessRepository:              PastIllnessRepository,
		PastInjuryRepository:               PastInjuryRepository,
		PastOptSurgeryRepository:           PastOptSurgeryRepository,
		PastSurgeryRepository:              PastSurgeryRepository,
		PatientChartRepository:             PatientChartRepository,
		PatientDiagnosisRepository:         PatientDiagnosisRepository,
		PatientEncounterLimitRepository:    PatientEncounterLimitRepository,
		PatientHistoryRepository:           PatientHistoryRepository,
		PatientQueueRepository:             PatientQueueRepository,
		PatientRepository:                  PatientRepository,
		PaymentWaiverRepository:            PaymentWaiverRepository,
		PaymentRepository:                  PaymentRepository,
		PharmacyRepository:                 PharmacyRepository,
		PhysicalExamFindingRepository:      PhysicalExamFindingRepository,
		PupilsRepository:                   PupilsRepository,
		QueueDestinationRepository:         QueueDestinationRepository,
		QueueSubscriptionRepository:        QueueSubscriptionRepository,
		ReferralOrderRepository:            ReferralOrderRepository,
		ReferralRepository:                 ReferralRepository,
		ReviewOfSystemRepository:           ReviewOfSystemRepository,
		RoomRepository:                     RoomRepository,
		SlitLampExamRepository:             SlitLampExamRepository,
		SupplyRepository:                   SupplyRepository,
		SurgicalOrderRepository:            SurgicalOrderRepository,
		SurgicalProcedureTypeRepository:    SurgicalProcedureTypeRepository,
		SurgicalProcedureRepository:        SurgicalProcedureRepository,
		SystemSymptomRepository:            SystemSymptomRepository,
		SystemRepository:                   SystemRepository,
		TreatmentOrderRepository:           TreatmentOrderRepository,
		TreatmentTypeRepository:            TreatmentTypeRepository,
		TreatmentRepository:                TreatmentRepository,
		UserTypeRepository:                 UserTypeRepository,
		UserRepository:                     UserRepository,
		VisitTypeRepository:                VisitTypeRepository,
		VisualAcuityRepository:             VisualAcuityRepository,
		VitalSignsRepository:               VitalSignsRepository,
	}}))

	r := gin.Default()
	//r.Use(cors.Default())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.GinContextToContextMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	authApi := auth.AuthApi{UserRepository: UserRepository}
	patientQueueApi := controller.PatientQueueApi{PatientQueueRepository: PatientQueueRepository, AppointmentRepository: AppointmentRepository}
	userTypeApi := controller.UserTypeApi{UserTypeRepository: UserTypeRepository}
	organizationDetailsApi := controller.OrganizationDetailsApi{OrganizationDetailsRepository: OrganizationDetailsRepository}

	r.Group("/public")
	{
		r.POST("/login", authApi.Login())
		r.POST("/legacy-login", authApi.LegacyLogin())
		r.POST("/signup", authApi.Signup)
		r.GET("/userTypes", userTypeApi.GetUserTypes)
		r.GET("/patientQueues", patientQueueApi.GetPatientQueues)
		r.GET("/organizationDetails", organizationDetailsApi.GetOrganizationDetails)

		r.Static("/files", "./files")

		r.GET("/rxnorm-drugs", controller.GetDrugs)
		r.GET("/rxnorm-intractions", controller.GetDrugIntractions)
	}

	r.Use(middleware.AuthMiddleware())
	r.GET("/api", playgroundHandler())
	r.POST("/query", graphqlHandler(s, h))

	return r
}

// GetDB returns gorm (ORM)
func (s *Server) GetDB() *gorm.DB {
	return s.DB
}

// GetConfig return the current app configuration
func (s *Server) GetConfig() *conf.Configuration {
	return s.Config
}

// GetModelRegistry returns the model registry
func (s *Server) GetModelRegistry() *models.Model {
	return s.ModelRegistry
}

// Start the http server
func (s *Server) Start() error {
	port := os.Getenv("ADDRESS")

	log.Fatal(s.Gin.Run(":" + port))
	return nil
}

func (s *Server) NewEnforcer() error {
	var model string
	var policy string

	appMode := os.Getenv("APP_MODE")

	if appMode == "release" {
		model = "/model.conf"
		policy = "/policy.csv"
	} else {
		model = "pkg/conf/model.conf"
		policy = "pkg/conf/policy.csv"
	}

	e, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		log.Fatal(err)
	}

	s.ACLEnforcer = e
	return nil
}

// GracefulShutdown Wait for interrupt signal
// to gracefully shutdown the server with a timeout of 5 seconds.
func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close database connection
	if s.DB != nil {
		db, _ := s.DB.DB()
		db.Close()
	}
}

// ShutdownTest shuts down test server
func (s *Server) ShutdownTest() {
	// close database connection
	if s.TestDB != nil {
		s.ModelRegistry.DropAll()
		db, _ := s.TestDB.DB()
		db.Close()
	}
}

// SeedData ...
func (m *Server) SeedData() error {
	appointmentStatusRepository := repository.ProvideAppointmentStatusRepository(m.DB)
	appointmentStatusRepository.Seed()

	userTypeRepository := repository.ProvideUserTypeRepository(m.DB)
	userTypeRepository.Seed()

	userRepository := repository.ProvideUserRepository(m.DB, userTypeRepository)
	userRepository.Seed()

	visitTypeRepository := repository.ProvideVisitTypeRepository(m.DB)
	visitTypeRepository.Seed()

	roomRepository := repository.ProvideRoomRepository(m.DB)
	roomRepository.Seed()

	billingRepository := repository.ProvideBillingRepository(m.DB)
	billingRepository.Seed()

	patientQueueRepository := repository.ProvidePatientQueueRepository(m.DB)
	patientQueueRepository.Seed()

	return nil
}

// AddSearchIndex ...
func (s *Server) AddSearchIndex() error {
	m := gormigrate.New(s.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// add search index
		{
			ID: "2021102815",
			Migrate: func(d *gorm.DB) error {
				// =======================
				// patients search index
				// =======================

				if err := d.Exec("UPDATE patients SET document = to_tsvector(id || ' ' || first_name || ' ' || last_name || ' ' || phone_no || ' ' || phone_no2 || ' ' || coalesce(email, ''))").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX patients_document_idx ON patients USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION patients_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.id || ' ' || new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.phone_no2 || ' ' || coalesce(new.email, ''));
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON patients FOR EACH ROW EXECUTE PROCEDURE patients_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// users search index
				// =======================

				if err := d.Exec("UPDATE users SET document = to_tsvector(id || ' ' || first_name || ' ' || last_name || ' ' || coalesce(email, ''))").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX users_document_idx ON users USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION users_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.id || ' ' || new.first_name || ' ' || new.last_name || ' ' || coalesce(new.email, ''));
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE users_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// billings search index
				// =======================

				if err := d.Exec("UPDATE billings SET document = to_tsvector(id || ' ' || item || ' ' || code)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX billings_document_idx ON billings USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION billings_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.id || ' ' || new.item || ' ' || new.code);
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON billings FOR EACH ROW EXECUTE PROCEDURE billings_tsvector_trigger()").Error; err != nil {

				}

				// ============================
				// appointments search index
				// ============================

				if err := d.Exec("UPDATE appointments SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX appointments_document_idx ON appointments USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION appointments_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no);
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON appointments FOR EACH ROW EXECUTE PROCEDURE appointments_tsvector_trigger()").Error; err != nil {

				}

				// ========================================
				// diagnostic procedure order search index
				// ========================================

				if err := d.Exec("UPDATE diagnostic_procedure_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX diagnostic_procedure_orders_document_idx ON diagnostic_procedure_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION diagnostic_procedure_orders_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON diagnostic_procedure_orders FOR EACH ROW EXECUTE PROCEDURE diagnostic_procedure_orders_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// lab order search index
				// =======================
				if err := d.Exec("UPDATE lab_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX lab_orders_document_idx ON lab_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION lab_orders_tsvector_trigger() RETURNS trigger AS $$
					begin
						new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
						return new;
					end
					$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON lab_orders FOR EACH ROW EXECUTE PROCEDURE lab_orders_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// diagnosis search index
				// =======================
				if err := d.Exec("UPDATE diagnoses SET document = to_tsvector(id || ' ' || coalesce(category_code, '') || ' ' || coalesce(diagnosis_code, '') || ' ' || coalesce(full_code, '') || ' ' || coalesce(full_description, '') || ' ' || coalesce(category_title, ''))").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX diagnoses_document_idx ON diagnoses USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION diagnoses_tsvector_trigger() RETURNS trigger AS $$
				begin
					new.document := to_tsvector(new.id || ' ' || coalesce(new.category_code, '') || ' ' || coalesce(new.diagnosis_code, '') || ' ' || coalesce(new.full_code, '') || ' ' || coalesce(new.full_description, '') || ' ' || coalesce(new.category_title, ''));
					return new;
				end
				$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON diagnoses FOR EACH ROW EXECUTE PROCEDURE diagnoses_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// treatment order search index
				// =======================

				if err := d.Exec("UPDATE treatment_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX treatment_orders_document_idx ON treatment_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION treatment_orders_tsvector_trigger() RETURNS trigger AS $$
					begin
						new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
						return new;
					end
					$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON treatment_orders FOR EACH ROW EXECUTE PROCEDURE treatment_orders_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// surgical order search index
				// =======================

				if err := d.Exec("UPDATE surgical_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX surgical_orders_document_idx ON surgical_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION surgical_orders_tsvector_trigger() RETURNS trigger AS $$
					begin
						new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
						return new;
					end
					$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON surgical_orders FOR EACH ROW EXECUTE PROCEDURE surgical_orders_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// follow up search index
				// =======================

				if err := d.Exec("UPDATE follow_up_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX follow_up_orders_document_idx ON follow_up_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION follow_up_orders_tsvector_trigger() RETURNS trigger AS $$
					begin
						new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
						return new;
					end
					$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON follow_up_orders FOR EACH ROW EXECUTE PROCEDURE follow_up_orders_tsvector_trigger()").Error; err != nil {

				}

				// =======================
				// referral search index
				// =======================
				if err := d.Exec("UPDATE referral_orders SET document = to_tsvector(first_name || ' ' || last_name || ' ' || phone_no || ' ' || user_name)").Error; err != nil {

				}

				if err := d.Exec("CREATE INDEX referral_orders_document_idx ON referral_orders USING GIN (document)").Error; err != nil {

				}

				if err := d.Exec(`CREATE FUNCTION referral_orders_tsvector_trigger() RETURNS trigger AS $$
					begin
						new.document := to_tsvector(new.first_name || ' ' || new.last_name || ' ' || new.phone_no || ' ' || new.user_name);
						return new;
					end
					$$ LANGUAGE plpgsql`).Error; err != nil {

				}

				if err := d.Exec("CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON referral_orders FOR EACH ROW EXECUTE PROCEDURE referral_orders_tsvector_trigger()").Error; err != nil {

				}

				return nil
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropColumn("patients", "document")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	return nil
}
