//visualization: netlistsvg

`include "comparator.sv"

module alu_comparator_half (
    input  logic [3:0] a, b,
    input  logic       prev_eq, prev_greater,
    output logic       eq, greater,
);

    logic internal_eq, internal_greater;
    comparator #(4) comp(
        .a(a), .b(b),
        .prev_eq(prev_eq),
        .prev_greater(prev_greater),
        .eq(internal_eq), .greater(internal_greater)
    );

    assign eq = internal_eq;
    assign greater = internal_greater;

endmodule