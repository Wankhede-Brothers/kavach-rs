pub mod bank;
pub mod write;
pub mod rpc;
pub mod stm;
pub mod kanban;
pub mod sync;
pub mod view;
pub mod inject;
pub mod spec;

use clap::Subcommand;

#[derive(Subcommand)]
pub enum MemoryCommand {
    Bank(bank::BankArgs),
    Write(write::WriteArgs),
    /// JSON-RPC on stdin
    Rpc,
    Stm,
    Inject,
    Spec,
    Kanban(kanban::KanbanArgs),
    View,
    Sync(sync::SyncArgs),
}

pub fn dispatch(cmd: MemoryCommand) -> anyhow::Result<()> {
    match cmd {
        MemoryCommand::Bank(args) => bank::run(args),
        MemoryCommand::Write(args) => write::run(args),
        MemoryCommand::Rpc => rpc::run(),
        MemoryCommand::Stm => stm::run(),
        MemoryCommand::Inject => inject::run(),
        MemoryCommand::Spec => spec::run(),
        MemoryCommand::Kanban(args) => kanban::run(args),
        MemoryCommand::View => view::run(),
        MemoryCommand::Sync(args) => sync::run(args),
    }
}
