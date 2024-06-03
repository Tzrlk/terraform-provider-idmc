package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &V2ActivityLogDataSource{}

func NewV2ActivityLogDataSource() datasource.DataSource {
	return &V2ActivityLogDataSource{}
}

// V2ActivityLogDataSource defines the data source implementation.
type V2ActivityLogDataSource struct {
	client *http.Client
}

type V2ActivityLogEntryDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Type                 types.String `tfsdk:"type"`
	ObjectId             types.String `tfsdk:"object_id"`
	ObjectName           types.String `tfsdk:"object_name"`
	RunId                types.Int64  `tfsdk:"run_id"`
	AgentId              types.String `tfsdk:"agent_id"`
	RuntimeEnvironmentId types.String `tfsdk:"runtime_environment_id"`
	StartTime            types.String `tfsdk:"start_time"`
	EndTime              types.String `tfsdk:"end_time"`
	StartTimeUtc         types.String `tfsdk:"start_time_utc"`
	EndTimeUtc           types.String `tfsdk:"end_time_utc"`
	State                types.Int64  `tfsdk:"state"`
	UIState              types.Int64  `tfsdk:"ui_state"` // only for specific id
	IsStopped            types.Bool   `tfsdk:"is_stopped"` // only for specific id
	FailedSourceRows     types.Int64  `tfsdk:"failed_source_rows"`
	SuccessSourceRows    types.Int64  `tfsdk:"success_source_rows"`
	FailedTargetRows     types.Int64  `tfsdk:"failed_target_rows"`
	SuccessTargetRows    types.Int64  `tfsdk:"success_target_rows"`
	ErrorMsg             types.String `tfsdk:"error_msg"` // only for specific id
	StartedBy            types.String `tfsdk:"started_by"` // only for specific id
	RunContextType       types.String `tfsdk:"run_context_type"` // only for specific id
	ScheduleName         types.String `tfsdk:"schedule_name"`
	OrgId                types.String `tfsdk:"org_id"` // only for specific id
	TotalSuccessRows     types.Int64  `tfsdk:"total_success_rows"` // only for specific id
	TotalFailedRows      types.Int64  `tfsdk:"total_failed_rows"` // only for specific id
	LogFilename          types.String `tfsdk:"log_filename"` // only for specific id
	ErrorFilename        types.String `tfsdk:"error_filename"` // only for specific id
	ErrorFileDir         types.String `tfsdk:"error_file_dir"` // only for specific id
	ConnType             types.String `tfsdk:"conn_type"` // only for specific id
	StopOnError          types.Bool   `tfsdk:"stop_on_error"` // only for specific id
	//Entries
}

// V2ActivityLogDataSourceModel describes the data source data model.
type V2ActivityLogDataSourceModel struct {
	Id       types.String                      `tfsdk:"id"`
	LogId    types.String                      `tfsdk:"log_id"`
	RunId    types.String                      `tfsdk:"run_id"`
	TaskId   types.String                      `tfsdk:"task_id"`
	Offset   types.Int64                       `tfsdk:"offset"`
	RowLimit types.Int64                       `tfsdk:"row_limit"`
	Session  types.Bool                        `tfsdk:"session"`
	Entries  V2ActivityLogEntryDataSourceModel `tfsdk:"entries"`
}

func (d *V2ActivityLogDataSource) Metadata(
		ctx context.Context,
		req datasource.MetadataRequest,
		resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v2_activity_log"
}

func (d *V2ActivityLogDataSource) Schema(
		ctx context.Context,
		req datasource.SchemaRequest,
		resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-2-resources/activitylog.html",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "V2ActivityLog identifier",
				Computed:            true,
			},
			"log_id": schema.StringAttribute{
				MarkdownDescription: "Log entry ID.\nInclude this attribute if you want to receive information for a specific ID.",
				Optional:            true,
			},
			"run_id": schema.StringAttribute{
				MarkdownDescription: "Job ID associated with the log entry ID.",
				Optional:            true,
			},
			"task_id": schema.StringAttribute{
				MarkdownDescription: "Task ID associated with the log entry ID. If taskId is not specified, all activityLog entries for all tasks are returned.",
				Optional:            true,
			},
			"offset": schema.StringAttribute{
				MarkdownDescription: "The number of rows to skip. For example, you might want to skip the first three rows.",
				Optional:            true,
			},
			"row_limit": schema.StringAttribute{
				MarkdownDescription: "The maximum number of rows to return. The maximum number you can specify is 1000.\nIf you omit this attribute, the activityLog returns all available rows, up to a maximum of 200 rows.",
				Optional:            true,
			},
			"session": schema.StringAttribute{
				MarkdownDescription: "True if a specific session log should be fetched. Requires log_id.",
				Optional:            true,
			},
			//"entries": schema.???{
			//	MarkdownDescription: "",
			//	Computed:            true,
			//},
		},
	}
}

func (d *V2ActivityLogDataSource) Configure(
		ctx context.Context,
		req datasource.ConfigureRequest,
		resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *V2ActivityLogDataSource) Read(
		ctx context.Context,
		req datasource.ReadRequest,
		resp *datasource.ReadResponse) {
	var data V2ActivityLogDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	//httpResp, err := d.client.Do(httpReq)
	//if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read v2ActivityLog, got error: %s", err))
	//	return
	//}
	//
	//httpResp.Body

	// For the purposes of this v2ActivityLog code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("v2ActivityLog-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
