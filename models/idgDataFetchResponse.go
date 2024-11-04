package models

type IdgDataFetchResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		EqualDecision       string `json:"equalDecision"`
		Timestamp           int64  `json:"timestamp"`
		Version             string `json:"version"`
		TransactionID       string `json:"transaction_id"`
		ConsumerInformation struct {
			ConsumerID string `json:"consumer_id"`
			Mobile     string `json:"mobile"`
		} `json:"consumer_information"`
		IdempotencyID    string `json:"idempotency_id"`
		ResponseMetadata struct {
			PartnerID string `json:"partnerId"`
		} `json:"response_metadata"`
		IDVerificationReport string `json:"id_verification_report"`
		KeyDetails           struct {
			BANKSTATEMENT struct {
				KeyName   string `json:"key_name"`
				KeyStatus string `json:"key_status"`
				KeyData   []struct {
					AccountNumber          string  `json:"account_number"`
					AccountType            string  `json:"account_type"`
					KeyID                  string  `json:"key_id"`
					KeyFetchType           string  `json:"key_fetch_type"`
					KeyVerificationStage   string  `json:"key_verification_stage"`
					VerificationType       string  `json:"verification_type"`
					MonthWiseSalaries      string  `json:"month_wise_salaries"`
					TransactionList        string  `json:"transaction_list"`
					IssuerName             string  `json:"issuer_name"`
					KeyName                string  `json:"key_name"`
					BankStatementEndDate   string  `json:"bank_statement_end_date"`
					BankStatementStartDate string  `json:"bank_statement_start_date"`
					ClosingBalance         float64     `json:"closing_balance"`
					BankName               string  `json:"bank_name"`
					Name                   string  `json:"name"`
					KeyBuid                string  `json:"key_buid"`
					KeySource              string  `json:"key_source"`
					OpeningBalance         float64 `json:"opening_balance"`
					MonthWiseAnalysis      []struct {
						MonthName              string  `json:"month_name"`
						NoOfDebitTransactions  int     `json:"no_of_debit_transactions"`
						NoOfCreditTransactions int     `json:"no_of_credit_transactions"`
						TotalCreditAmount      float64     `json:"total_credit_amount"`
						Year                   string  `json:"year"`
						TotalDebitAmount       float64 `json:"total_debit_amount"`
						AverageEodBalance      float64     `json:"average_eod_balance"`
					} `json:"month_wise_analysis"`
					KeyFetchedAt string `json:"key_fetched_at"`
				} `json:"key_data"`
			} `json:"BANK_STATEMENT"`
			PAN struct {
				KeyName   string `json:"key_name"`
				KeyStatus string `json:"key_status"`
				KeyData   []struct {
					KeyName              string `json:"key_name"`
					Gender               string `json:"gender"`
					KeyID                string `json:"key_id"`
					Dob                  string `json:"dob"`
					Name                 string `json:"name"`
					KeyBuid              string `json:"key_buid"`
					KeyVerificationStage string `json:"key_verification_stage"`
					KeySource            string `json:"key_source"`
					KeyFetchedAt         string `json:"key_fetched_at"`
				} `json:"key_data"`
			} `json:"PAN"`
		} `json:"key_details"`
	} `json:"data"`
	StatusCode string `json:"status_code"`
}
