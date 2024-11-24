package canparser

import (
	"fmt"
	"strconv"
	"strings"
)

// CANMessage представляет одно сообщение CAN
type CANMessage struct {
	Interface string `json:"interface"`
	ID        uint32 `json:"id"`
	Length    int    `json:"length"`
	Data      []byte `json:"data"`
	EDP       int    `json:"edp"`
	DP        int    `json:"dp"`
	PF        int    `json:"pf"`
	PS        int    `json:"ps"`
	SA        int    `json:"sa"`
	PGN       int    `json:"pgn"`
	Raw       string `json:"raw"`
}

// Parser представляет парсер для сообщений CAN
type Parser struct {
}

// NewParser создает новый экземпляр парсера
func NewParser() *Parser {
	return &Parser{}
}

// ParseLine парсит одну строку из дампа CAN в структуру CANMessage
func (p *Parser) ParseLine(line string) (*CANMessage, error) {
	// Разделение строки на части
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	// Интерфейс
	iface := parts[0]

	// ID сообщения
	idHex := parts[1]
	id, err := strconv.ParseUint(idHex, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid CAN ID: %s", idHex)
	}

	// Длина сообщения
	length := 0
	fmt.Sscanf(parts[2], "[%d]", &length)

	// Полезная нагрузка
	data := []byte{}
	for _, byteStr := range parts[3:] {
		b, err := strconv.ParseUint(byteStr, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid data byte: %s", byteStr)
		}
		data = append(data, byte(b))
	}

	// Расшифровка полей J1939
	edp := int((id >> 25) & 0x01)
	dp := int((id >> 24) & 0x01)
	pf := int((id >> 16) & 0xFF)
	ps := int((id >> 8) & 0xFF)
	sa := int(id & 0xFF)

	// Вычисление PGN
	var pgn int
	if pf >= 240 {
		pgn = pf << 8 // глобальный адрес
	} else {
		pgn = (pf << 8) | ps // адресный диапазон
	}

	return &CANMessage{
		Interface: iface,
		ID:        uint32(id),
		Length:    length,
		Data:      data,
		EDP:       edp,
		DP:        dp,
		PF:        pf,
		PS:        ps,
		SA:        sa,
		PGN:       pgn,
		Raw:       line,
	}, nil
}
