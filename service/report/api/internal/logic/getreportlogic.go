// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"os"
	"time"

	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/types"
	"github.com/Drengin6306/ZeroBank/service/transaction/rpc/transaction"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 生成报表
func NewGetReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReportLogic {
	return &GetReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReportLogic) GetReport(req *types.GenerateReportRequest) (resp *types.GenerateReportResponse, err error) {
	// 获取用户
	account := l.ctx.Value(vars.AccountKey).(string)

	// 获取交易流水
	transactions, err := l.svcCtx.TransactionRpc.GetTransactions(l.ctx, &transaction.TransactionRequest{
		AccountId: account,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}
	// 生成excel报表
	fileName, err := generateExcel(transactions.Transactions)
	if err != nil {
		return nil, err
	}
	return &types.GenerateReportResponse{
		FileName: fileName,
	}, nil
}

func generateExcel(records []*transaction.Record) (name string, err error) {
	// 取款 存款 收到他人转账 向他人转账
	f := excelize.NewFile()
	defer f.Close()

	// 1. 创建流式写入器 (针对 Sheet1)
	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return "", err
	}

	// 2. 定义表头样式 (加粗，背景色)
	styleID, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
	})
	// 设置列宽
	if err := sw.SetColWidth(1, 1, 22); err != nil { // A列：交易流水号
		return "", err
	}
	if err := sw.SetColWidth(2, 2, 22); err != nil { // B列：付款账户
		return "", err
	}
	if err := sw.SetColWidth(3, 3, 22); err != nil { // C列：收款账户
		return "", err
	}
	if err := sw.SetColWidth(4, 4, 12); err != nil { // D列：交易类型
		return "", err
	}
	if err := sw.SetColWidth(5, 5, 12); err != nil { // E列：交易金额
		return "", err
	}
	if err := sw.SetColWidth(6, 6, 25); err != nil { // F列：交易时间
		return "", err
	}

	// 3. 写入表头
	headers := []interface{}{
		excelize.Cell{Value: "交易流水号", StyleID: styleID},
		excelize.Cell{Value: "付款账户", StyleID: styleID},
		excelize.Cell{Value: "收款账户", StyleID: styleID},
		excelize.Cell{Value: "交易类型", StyleID: styleID},
		excelize.Cell{Value: "交易金额", StyleID: styleID},
		excelize.Cell{Value: "交易时间", StyleID: styleID},
	}
	if err := sw.SetRow("A1", headers); err != nil {
		return "", err
	}

	// 4. 遍历数据并写入
	for i, record := range records {
		rowNum := i + 2
		axis, _ := excelize.CoordinatesToCellName(1, rowNum) // 列从1开始，行从2开始 A2, A3, ...

		// 处理空值逻辑
		accTo := *record.AccountTo
		if accTo == "" {
			accTo = "-"
		}

		row := []interface{}{
			record.TransactionId,
			record.AccountFrom,
			accTo,
			typeToString(record.TransactionType), // 转换枚举值
			record.Amount,                        // 金额
			formatTime(record.CreatedAt),         // 格式化时间
		}

		if err := sw.SetRow(axis, row); err != nil {
			return "", err
		}
	}

	// 5. 刷新流，完成写入
	if err := sw.Flush(); err != nil {
		return "", err
	}

	// 6. 保存文件到temp目录
	fileName := "transaction_report_" + time.Now().Format("20060102150405") + ".xlsx"
	// 没有的话创建temp目录
	if err := os.MkdirAll("./temp", 0755); err != nil {
		return "", err
	}
	if err := f.SaveAs("./temp/" + fileName); err != nil {
		return "", err
	}

	return fileName, nil
}

func typeToString(t int32) string {
	switch t {
	case vars.TransactionTypeDeposit:
		return "存款"
	case vars.TransactionTypeWithdraw:
		return "取款"
	case vars.TransactionTypeTransfer:
		return "转账"
	default:
		return "Unknown"
	}
}

func formatTime(ts string) string {
	// RFC3339 格式的时间字符串
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts // 如果解析失败，返回原始字符串
	}
	return parsedTime.Format("2006-01-02 15:04:05")
}
