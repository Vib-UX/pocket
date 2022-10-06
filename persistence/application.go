package persistence

import (
	"encoding/hex"
	"log"

	"github.com/pokt-network/pocket/persistence/types"
	"google.golang.org/protobuf/proto"

	"github.com/pokt-network/pocket/shared/modules"
)

func (p PostgresContext) UpdateApplicationsTree(apps []modules.Actor) error {
	for _, app := range apps {
		bzAddr, err := hex.DecodeString(app.GetAddress())
		if err != nil {
			return err
		}

		appBz, err := proto.Marshal(app.(*types.Actor))
		if err != nil {
			return err
		}

		// OPTIMIZE: This is the only line unique to `Application`
		if _, err := p.MerkleTrees[appMerkleTree].Update(bzAddr, appBz); err != nil {
			return err
		}
	}

	return nil
}

func (p PostgresContext) getApplicationsUpdatedAtHeight(height int64) (apps []*types.Actor, err error) {
	// OPTIMIZE: This is the only line unique to `Application`
	actors, err := p.GetActorsUpdated(types.ApplicationActor, height)
	if err != nil {
		return nil, err
	}

	apps = make([]*types.Actor, len(actors))
	for _, actor := range actors {
		app := &types.Actor{
			ActorType:       types.ActorType_App,
			Address:         actor.Address,
			PublicKey:       actor.PublicKey,
			Chains:          actor.Chains,
			GenericParam:    actor.ActorSpecificParam,
			StakedAmount:    actor.StakedTokens,
			PausedHeight:    actor.PausedHeight,
			UnstakingHeight: actor.UnstakingHeight,
			Output:          actor.OutputAddress,
		}
		apps = append(apps, app)
	}
	return
}

func (p PostgresContext) GetAppExists(address []byte, height int64) (exists bool, err error) {
	return p.GetExists(types.ApplicationActor, address, height)
}

func (p PostgresContext) GetApp(address []byte, height int64) (operator, publicKey, stakedTokens, maxRelays, outputAddress string, pauseHeight, unstakingHeight int64, chains []string, err error) {
	actor, err := p.GetActor(types.ApplicationActor, address, height)
	if err != nil {
		return
	}
	operator = actor.Address
	publicKey = actor.PublicKey
	stakedTokens = actor.StakedTokens
	maxRelays = actor.ActorSpecificParam
	outputAddress = actor.OutputAddress
	pauseHeight = actor.PausedHeight
	unstakingHeight = actor.UnstakingHeight
	chains = actor.Chains
	return
}

func (p PostgresContext) InsertApp(address []byte, publicKey []byte, output []byte, _ bool, _ int32, maxRelays string, stakedTokens string, chains []string, pausedHeight int64, unstakingHeight int64) error {
	return p.InsertActor(types.ApplicationActor, types.BaseActor{
		Address:            hex.EncodeToString(address),
		PublicKey:          hex.EncodeToString(publicKey),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: maxRelays,
		OutputAddress:      hex.EncodeToString(output),
		PausedHeight:       pausedHeight,
		UnstakingHeight:    unstakingHeight,
		Chains:             chains,
	})
}

func (p PostgresContext) UpdateApp(address []byte, maxRelays string, stakedAmount string, chains []string) error {
	return p.UpdateActor(types.ApplicationActor, types.BaseActor{
		Address:            hex.EncodeToString(address),
		StakedTokens:       stakedAmount,
		ActorSpecificParam: maxRelays,
		Chains:             chains,
	})
}

func (p PostgresContext) GetAppStakeAmount(height int64, address []byte) (string, error) {
	return p.GetActorStakeAmount(types.ApplicationActor, address, height)
}

func (p PostgresContext) SetAppStakeAmount(address []byte, stakeAmount string) error {
	return p.SetActorStakeAmount(types.ApplicationActor, address, stakeAmount)
}

func (p PostgresContext) DeleteApp(_ []byte) error {
	log.Println("[DEBUG] DeleteApp is a NOOP")
	return nil
}

func (p PostgresContext) GetAppsReadyToUnstake(height int64, _ int32) ([]modules.IUnstakingActor, error) {
	return p.GetActorsReadyToUnstake(types.ApplicationActor, height)
}

func (p PostgresContext) GetAppStatus(address []byte, height int64) (int32, error) {
	return p.GetActorStatus(types.ApplicationActor, address, height)
}

func (p PostgresContext) SetAppUnstakingHeightAndStatus(address []byte, unstakingHeight int64, status int32) error {
	return p.SetActorUnstakingHeightAndStatus(types.ApplicationActor, address, unstakingHeight)
}

func (p PostgresContext) GetAppPauseHeightIfExists(address []byte, height int64) (int64, error) {
	return p.GetActorPauseHeightIfExists(types.ApplicationActor, address, height)
}

func (p PostgresContext) SetAppStatusAndUnstakingHeightIfPausedBefore(pausedBeforeHeight, unstakingHeight int64, status int32) error {
	return p.SetActorStatusAndUnstakingHeightIfPausedBefore(types.ApplicationActor, pausedBeforeHeight, unstakingHeight)
}

func (p PostgresContext) SetAppPauseHeight(address []byte, height int64) error {
	return p.SetActorPauseHeight(types.ApplicationActor, address, height)
}

func (p PostgresContext) GetAppOutputAddress(operator []byte, height int64) ([]byte, error) {
	return p.GetActorOutputAddress(types.ApplicationActor, operator, height)
}
