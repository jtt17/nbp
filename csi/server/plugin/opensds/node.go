package opensds

import (
	"log"
	"os"
	"runtime"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/opensds/nbp/client/iscsi"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
//                            Node Service                                    //
////////////////////////////////////////////////////////////////////////////////

// NodePublishVolume implementation
func (p *Plugin) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	log.Println("start to NodePublishVolume")
	defer log.Println("end to NodePublishVolume")

	portal := req.PublishVolumeInfo["portal"]
	targetiqn := req.PublishVolumeInfo["targetiqn"]
	targetlun := req.VolumeId

	// Connect Target
	device, err := iscsi.Connect(portal, targetiqn, targetlun)
	if err != nil {
		return nil, err
	}

	// Format and Mount
	err = iscsi.FormatandMount(device, "", req.TargetPath)
	if err != nil {
		return nil, err
	}

	return &csi.NodePublishVolumeResponse{
		Reply: &csi.NodePublishVolumeResponse_Result_{
			Result: &csi.NodePublishVolumeResponse_Result{},
		},
	}, nil
}

// NodeUnpublishVolume implementation
func (p *Plugin) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	log.Println("start to NodeUnpublishVolume")
	defer log.Println("end to NodeUnpublishVolume")

	// Umount
	err := iscsi.Umount(req.TargetPath)
	if err != nil {
		return nil, err
	}

	// Disconnect
	// TODO: get portal and targetiqn
	err = iscsi.Disconnect("", "")
	if err != nil {
		return nil, err
	}

	return &csi.NodeUnpublishVolumeResponse{
		Reply: &csi.NodeUnpublishVolumeResponse_Result_{
			Result: &csi.NodeUnpublishVolumeResponse_Result{},
		},
	}, nil
}

// GetNodeID implementation
func (p *Plugin) GetNodeID(
	ctx context.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {

	log.Println("start to GetNodeID")
	defer log.Println("end to GetNodeID")

	// Get host name from os
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &csi.GetNodeIDResponse{
		Reply: &csi.GetNodeIDResponse_Result_{
			Result: &csi.GetNodeIDResponse_Result{
				NodeId: hostname,
			},
		},
	}, nil
}

// NodeProbe implementation
func (p *Plugin) NodeProbe(
	ctx context.Context,
	req *csi.NodeProbeRequest) (
	*csi.NodeProbeResponse, error) {

	log.Println("start to NodeProbe")
	defer log.Println("end to NodeProbe")

	switch runtime.GOOS {
	case "linux":
		return &csi.NodeProbeResponse{
			Reply: &csi.NodeProbeResponse_Result_{
				Result: &csi.NodeProbeResponse_Result{},
			},
		}, nil
	default:
		msg := "unsupported operating system:" + runtime.GOOS
		log.Fatalf(msg)
		return &csi.NodeProbeResponse{
			Reply: &csi.NodeProbeResponse_Error{
				Error: &csi.Error{
					Value: &csi.Error_NodeProbeError_{
						NodeProbeError: &csi.Error_NodeProbeError{
							ErrorCode:        csi.Error_NodeProbeError_MISSING_REQUIRED_HOST_DEPENDENCY,
							ErrorDescription: msg,
						},
					},
				},
			},
		}, nil
	}
}

// NodeGetCapabilities implementation
func (p *Plugin) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	log.Println("start to NodeGetCapabilities")
	defer log.Println("end to NodeGetCapabilities")

	return &csi.NodeGetCapabilitiesResponse{
		Reply: &csi.NodeGetCapabilitiesResponse_Result_{
			Result: &csi.NodeGetCapabilitiesResponse_Result{
				Capabilities: []*csi.NodeServiceCapability{
					&csi.NodeServiceCapability{
						Type: &csi.NodeServiceCapability_Rpc{
							Rpc: &csi.NodeServiceCapability_RPC{
								Type: csi.NodeServiceCapability_RPC_UNKNOWN,
							},
						},
					},
				},
			},
		},
	}, nil
}
