package common

import (
	"encoding/json"
)

type tunnelInfoDto struct {
	Id                string `json:"id"`
	TunnelType        string `json:"tunnel_type"`
	EdgeServerAddress string `json:"edge_server_address"`
	EdgeAgentAddress  string `json:"edge_agent_address"`
}

type tunnelInfo struct {
	id         Id
	tunnelType TunnelType
	esAddr     string
	eaAddr     string
}

func GenerateTunnelInfo(tunnelId Id, tunnelType TunnelType, edgeServerAddr string, edgeAgentAddr string) (
	TunnelInfo, error,
) {
	if tunnelType != TunnelTypeTcp && tunnelType != TunnelTypeHttp && tunnelType != TunnelTypeUdp {
		return nil, ErrTunnelInfoInvalid
	}

	if edgeAgentAddr == "" || edgeServerAddr == "" {
		return nil, ErrTunnelInfoInvalid
	}

	if tunnelId == nil {
		tunnelId = GenerateId()
	}

	return &tunnelInfo{
		id:         tunnelId,
		tunnelType: tunnelType,
		esAddr:     edgeServerAddr,
		eaAddr:     edgeAgentAddr,
	}, nil
}

func MarshalTunnelInfo(tunnelInfo TunnelInfo) ([]byte, error) {
	if tunnelInfo == nil {
		return nil, ErrTunnelInfoInvalid
	}
	dto := &tunnelInfoDto{
		Id:                tunnelInfo.Id().String(),
		TunnelType:        string(tunnelInfo.Type()),
		EdgeServerAddress: tunnelInfo.EdgeServerAddr(),
		EdgeAgentAddress:  tunnelInfo.EdgeAgentAddr(),
	}
	data, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalTunnelInfo(data []byte) (TunnelInfo, error) {
	dto := &tunnelInfoDto{}
	err := json.Unmarshal(data, dto)
	if err != nil {
		return nil, err
	}
	tunnelId, err := ParseId(dto.Id)
	if err != nil {
		return nil, err
	}
	tunnelType := TunnelType(dto.TunnelType)
	ti, err := GenerateTunnelInfo(tunnelId, tunnelType, dto.EdgeServerAddress, dto.EdgeAgentAddress)
	if err != nil {
		return nil, err
	}
	return ti, nil
}

func (t *tunnelInfo) Id() Id {
	return t.id
}

func (t *tunnelInfo) Type() TunnelType {
	return t.tunnelType
}

func (t *tunnelInfo) EdgeServerAddr() string {
	return t.esAddr
}

func (t *tunnelInfo) EdgeAgentAddr() string {
	return t.eaAddr
}
