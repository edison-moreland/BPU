`include "inc_rst.sv"

module microop_counter #(
    N = 3
) (
    input logic clk, inst_done,
    output logic [N-1:0] count,
);

    logic [N-1:0] next_count;

    inc_rst #(N) ir(
        .rst(inst_done),
        .din(count), .dout(next_count)
    );

    always_ff @(negedge clk) begin
        count <= next_count;
    end
endmodule