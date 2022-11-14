use cardano_test

db.tx.drop()
db.utxo.drop()
db.stxi.drop()
db.asst.drop()
db.mint.drop()
db.meta.drop()
db.dtum.drop()

db.coll.drop()
db.rdmr.drop()
db.witn.drop()
db.trans.drop()
db.witp.drop()

db.blck.drop()
db.block.drop()

db.witp.drop()
db.witp.drop()

db.pool.drop()
db.dele.drop()
db.reti.drop()
db.skde.drop()
db.skre.drop()
db.cip25.drop()
db.cip15.drop()
db.scpt.drop()

db.cip25p.drop()


db.tx.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.tx.createIndex({"context.tx_hash":1})
db.tx.createIndex({"context.block_number":1})

db.utxo.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.utxo.createIndex({"context.tx_hash":1})
db.utxo.createIndex({"tx_output.address":1,"context.block_number":1})

db.stxi.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.stxi.createIndex({"context.tx_hash":1})

db.asst.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.asst.createIndex({"context.tx_hash":1})

db.mint.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.mint.createIndex({"context.tx_hash":1})
db.mint.createIndex({"mint.policy":1,"mint.asset":1})

db.meta.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.meta.createIndex({"context.tx_hash":1})
db.meta.createIndex({"metadata.label":1,"context.timestamp":1})
db.meta.createIndex({"metadata.map_json.msg":1}, { sparse: true } )


db.dtum.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.dtum.createIndex({"context.tx_hash":1})

db.blck.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.blck.createIndex({"context.block_number":1})

db.block.createIndex({"context.block_number":1,"context.block_hash":1},{unique:true})

db.blockf.createIndex({"height":1},{unique:true})
db.blockf.createIndex({"epoch":1})
db.blockf.createIndex({"slot_leader":1})

db.coll.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.coll.createIndex({"context.tx_hash":1})

db.rdmr.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1,"plutus_redeemer.purpose":1},{unique:true})
db.rdmr.createIndex({"context.tx_hash":1})

db.witn.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.witn.createIndex({"context.tx_hash":1})

db.witp.createIndex({"fingerprint":1,"context.block_hash":1,"context.tx_hash":1},{unique:true})
db.witp.createIndex({"context.tx_hash":1})



db.trans.createIndex({"fingerprint":1,"context.tx_hash":1,"context.block_hash":1},{unique:true})
db.trans.createIndex({"context.tx_hash":1})
db.trans.createIndex({"utxo_output.assets.policy":1,"utxo_output.assets.asset":1})
db.trans.createIndex({"context.block_number":1})
db.trans.createIndex({"utxo_output.address":1,"context.timestamp":-1})
db.trans.createIndex({"metadata.label":1,"context.timestamp":1})
db.trans.createIndex( {"metadata.map_json.msg": 1 }, { sparse: true } )
db.trans.createIndex( {"stake_delegation.pool_hash": 1 }, { sparse: true } )
db.trans.createIndex( {"stake_delegation.credential.AddrKeyhash": 1 }, { sparse: true } )


db.trans_nft.drop()
db.trans_nft.createIndex({"fingerprint":1,"context.tx_hash":1,"context.block_hash":1},{unique:true})
db.trans_nft.createIndex({"context.tx_hash":1})
db.trans_nft.createIndex({"context.block_number":1})
db.trans_nft.createIndex({"utxo_output.address":1,"context.timestamp":-1})
db.trans_nft.createIndex({"tx_meta.market_place.name":1})
db.trans_nft.createIndex({"utxo_output.assets.policy":1})
 
db.pool.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.pool.createIndex({"context.tx_hash":1})
db.pool.createIndex({"pool_registration.pool_id":1})

db.dele.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.dele.createIndex({"context.tx_hash":1})
db.dele.createIndex({"stake_delegation.pool_hash":1})


db.reti.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.reti.createIndex({"context.tx_hash":1})
db.reti.createIndex({"pool_retirement.pool":1})

db.skde.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.skde.createIndex({"context.tx_hash":1})
db.skde.createIndex({"stake_deregistration.credential.addrkey_hash":1})


db.skre.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.skre.createIndex({"context.tx_hash":1})
db.skre.createIndex({"stake_registration.credential.addrkey_hash":1})

db.cip25.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.cip25.createIndex({"context.tx_hash":1})
db.cip25.createIndex({"context.block_number":1})
db.cip25.createIndex({"cip25_asset.policy":1,"cip25_asset.asset":1})



db.cip25p.createIndex({"fingerprint":1,"context.tx_hash":1},{unique:true})
db.cip25p.createIndex({"context.tx_hash":1})
db.cip25p.createIndex({"context.block_number":1})
db.cip25p.createIndex({"cip25_asset.policy":1,"cip25_asset.asset":1})



db.cip15.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})
db.scpt.createIndex({"fingerprint":1,"context.block_hash":1},{unique:true})


